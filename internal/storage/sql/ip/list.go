package ip

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/domain/entities"
)

// For save into DB column of type BIT or BIT VARYING
type BitMask string

type Row struct {
	//IPBits      net.IPBits      `db:"ip"`
	//MaskBits    net.IPMask  `db:"mask"`

	IPBits   BitMask     `db:"ip"`
	MaskBits *BitMask    `db:"mask"` // pointer because could be nil (NULL)
	IP       entities.IP `db:"ip_str"`
	MaskInt  *int        `db:"mask_int"` // pointer because could be nil (NULL)
	IsV6     bool        `db:"is_v6"`
	Kind     string      `db:"kind"`
}

type List struct {
	db   *sqlx.DB
	kind string // to identify list, (for this program is "black" and "white")
}

func NewList(db *sqlx.DB, kind string) *List {
	return &List{
		db:   db,
		kind: kind,
	}
}

func (l *List) Add(ctx context.Context, ip entities.IP) error {
	// Check if this IP already in DB, it is forbid has 2 same IPs in DB, violet primary key constraint
	has, err := l.Has(ctx, ip)
	if err != nil {
		return wrapAddError(l.kind, err)
	}

	if has {
		return nil
	}

	query := `INSERT INTO ip_list(ip, mask, ip_str, mask_int, is_v6, kind) 
				VALUES(:ip, :mask, :ip_str, :mask_int, :is_v6, :kind)`

	row, err := convertIPToRow(ip, l.kind)

	if err != nil {
		return wrapAddError(l.kind, err)
	}

	_, err = l.db.NamedExecContext(ctx, query, row)
	if err != nil {
		return wrapAddError(l.kind, err)
	}

	return nil
}

func (l *List) Delete(ctx context.Context, ip entities.IP) error {
	row, err := convertIPToRow(ip, l.kind)
	if err != nil {
		return wrapDeleteError(l.kind, err)
	}

	condition := buildFindIPCondition(row)
	query := `DELETE FROM ip_list WHERE %s`
	query = fmt.Sprintf(query, condition)

	_, err = l.db.NamedExecContext(ctx, query, row)
	if err != nil {
		return wrapDeleteError(l.kind, err)
	}

	return nil
}

func (l *List) Has(ctx context.Context, ip entities.IP) (bool, error) {
	row, err := convertIPToRow(ip, l.kind)
	if err != nil {
		return false, wrapFindError(l.kind, err)
	}

	condition := buildFindIPCondition(row)

	query := `SELECT COUNT(*) FROM ip_list WHERE %s`
	query = fmt.Sprintf(query, condition)

	stmt, err := l.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return false, wrapFindError(l.kind, err)
	}

	var count int
	queryRow := stmt.QueryRowxContext(ctx, row)

	err = queryRow.Scan(&count)
	if err != nil {
		return false, wrapFindError(l.kind, err)
	}

	return count > 0, nil
}

func (l *List) IsConform(ctx context.Context, ip entities.IP) (bool, error) {
	if ip.HasMaskPart() {
		return false, wrapConformError(l.kind, ip, errors.New("expected pure IPBits (without mask)"))
	}

	// condition when search pure ip exact as it
	conditionFindExactIP := `ip = :ip AND mask_int IS NULL`

	// condition when find subnet (network) IPBits that contains this pure IPBits
	// see net.Contains method
	conditionFindNetworkContactsIP := `ip & mask = :ip & mask AND mask_int IS NOT NULL`

	// search in list of current kind
	conditionKind := `kind = :kind`

	// search among IPs of certain version
	conditionIPVersion := `is_v6 = :is_v6`

	// result query
	query := `SELECT ip FROM ip_list 
				WHERE ((%s) OR (%s)) AND %s AND %s
				LIMIT 1`

	query = fmt.Sprintf(query, conditionFindExactIP, conditionFindNetworkContactsIP, conditionKind, conditionIPVersion)

	row, err := convertIPToRow(ip, l.kind)
	if err != nil {
		return false, wrapConformError(l.kind, ip, err)
	}

	stmt, err := l.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return false, wrapConformError(l.kind, ip, err)
	}

	var resIP string
	queryRow := stmt.QueryRowxContext(ctx, row)

	err = queryRow.Scan(&resIP)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, wrapConformError(l.kind, ip, err)
	}

	return resIP != "", err
}

func (l *List) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM ip_list WHERE kind = $1`

	var count int
	queryRow := l.db.QueryRowContext(ctx, query, l.kind)

	err := queryRow.Scan(&count)
	if err != nil {
		return 0, wrapCountError(l.kind, err)
	}

	return count, nil
}

func (l *List) Clear(ctx context.Context) error {
	query := `DELETE FROM ip_list WHERE kind = $1`

	_, err := l.db.ExecContext(ctx, query, l.kind)
	if err != nil {
		return wrapClearError(l.kind, err)
	}

	return nil
}

func convertIPToRow(ip entities.IP, kind string) (Row, error) {
	row := Row{
		IP:   ip,
		Kind: kind,
	}

	var netIP net.IP

	if ip.HasMaskPart() {
		maskInt, _ := ip.GetMaskAsInt()
		row.MaskInt = &maskInt

		var err error
		var netIPNet *net.IPNet

		netIP, netIPNet, err = ip.ParseAsCIDR()
		if err != nil {
			return row, err
		}

		// mask of subnet
		mask := convertByteSliceToBitMask(netIPNet.Mask)
		row.MaskBits = &mask
	} else {
		netIP = ip.Parse()
		if netIP == nil {
			return row, fmt.Errorf("invalid ip `%s`", ip)
		}
	}

	// 4 byte representation
	b4IP := netIP.To4()

	// if nil then it is IPv6
	if b4IP != nil {
		row.IPBits = convertByteSliceToBitMask(b4IP)
		row.IsV6 = false
	} else {
		row.IPBits = convertByteSliceToBitMask(netIP)
		row.IsV6 = true
	}

	return row, nil
}

func convertByteSliceToBitMask(bs []byte) BitMask {
	lbs := len(bs)
	byteBits := make([]string, lbs)
	for i := 0; i < lbs; i++ {
		byteBits[i] = fmt.Sprintf("%08b", bs[i])
	}
	return BitMask(strings.Join(byteBits, ""))
}

func buildFindIPCondition(row Row) string {
	var condition string
	if row.MaskInt != nil {
		condition = `ip = :ip AND mask = :mask AND kind = :kind`
	} else {
		condition = `ip = :ip AND kind = :kind`
	}
	return condition
}
func wrapAddError(kind string, err error) error {
	return fmt.Errorf("failed to add IPBits in list %s: %w", kind, err)
}

func wrapDeleteError(kind string, err error) error {
	return fmt.Errorf("failed to delete IPBits from list %s: %w", kind, err)
}

func wrapFindError(kind string, err error) error {
	return fmt.Errorf("failed find IPBits in list %s: %w", kind, err)
}

func wrapCountError(kind string, err error) error {
	return fmt.Errorf("failed count IPs in list %s: %w", kind, err)
}

func wrapClearError(kind string, err error) error {
	return fmt.Errorf("failed to delete all IPs from list %s: %w", kind, err)
}

func wrapConformError(kind string, ip entities.IP, err error) error {
	return fmt.Errorf("failed when check `%s` is conform IPs in list %s: %w", ip, kind, err)
}
