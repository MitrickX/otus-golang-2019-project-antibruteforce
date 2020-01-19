package bucket

type Storage interface {
	Add(bucket Bucket, key interface{}) (bool, error)
	Delete(key interface{}) (bool, error)
	Get(key interface{}) (Bucket, error)
	Has(key interface{}) (bool, error)
	Count() (int, error)
}
