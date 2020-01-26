Feature: Auth
  As API client of anti brute force service
  I want check is auth allowed

  Scenario: Test auth when ip bucket is overflowing
    Given Clean bucket for
    """
    ip=127.0.0.1
    """
    And Clean "white" list
    And Clean "black" list

    # N auth tries with different random credentials, must be ok. N is limit for ip bucket
    When I call "ipLimit" times method "Auth" with params:
    """
    login=random&password=random&ip=127.0.0.1
    """
    Then The error must be "nil"
    And The result must be "true"

    # N + 1 auth try with same ip, must be not ok cause of overflowing
    When I call method "Auth" with params:
    """
    login=random&password=random&ip=127.0.0.1
    """
    Then The error must be "nil"
    And The result must be "false"

    # 1 auth try for different ip, must be ok cause new id - new bucket
    When I call method "Auth" with params:
    """
    login=random&password=random&ip=127.0.0.2
    """
    Then The error must be "nil"
    And The result must be "true"

  Scenario: Test auth when ip bucket is not overflowing
    Given Clean bucket for
    """
    ip=127.0.0.1
    """
    And Clean "white" list
    And Clean "black" list

    # N auth tries with different random credentials but with same ip must be ok. N is limit for ip bucket
    When I call "ipLimit" times method "Auth" with params:
    """
    login=random&password=random&ip=127.0.0.1
    """
    Then The error must be "nil"
    And The result must be "true"

    # wait 1 minute
    When Wait 1 minute

    # N auth tries with different random credentials but with same ip, must be ok. N is limit for ip bucket
    And I call "ipLimit" times method "Auth" with params:
    """
    login=random&password=random&ip=127.0.0.1
    """
    Then The error must be "nil"
    And The result must be "true"
