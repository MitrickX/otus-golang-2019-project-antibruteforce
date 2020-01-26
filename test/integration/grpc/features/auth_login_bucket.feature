Feature: Auth
  As API client of anti brute force service
  I want check is auth allowed

  Scenario: Test auth when login bucket is overflowing
    Given Clean bucket for
    """
    login=test&ip=127.0.0.1
    """
    And Clean "white" list
    And Clean "black" list

    # N auth tries, must be ok. N is limit for login bucket
    When I call "loginLimit" times method "Auth" with params:
    """
    login=test&ip=127.0.0.1
    """
    Then The error must be "nil"
    And The result must be "true"

    # N + 1 auth try with same login, must be not ok cause of overflowing
    When I call method "Auth" with params:
    """
    login=test&ip=127.0.0.1
    """
    Then The error must be "nil"
    And The result must be "false"

    # 1 auth try for different login, must be ok cause new login - new bucket
    When I call method "Auth" with params:
    """
    login=test2&ip=127.0.0.1
    """
    Then The error must be "nil"
    And The result must be "true"

  Scenario: Test auth when login bucket is not overflowing
    Given Clean bucket for
    """
    login=test&ip=127.0.0.1
    """
    And Clean "white" list
    And Clean "black" list

    # N auth tries, must be ok. N is limit for login bucket
    When I call "loginLimit" times method "Auth" with params:
    """
    login=test&ip=127.0.0.1
    """
    Then The error must be "nil"
    And The result must be "true"

    # "wait" 1 minute
    When Wait 1 minute

    # N auth tries, must be ok. N is limit for login bucket
    And I call "loginLimit" times method "Auth" with params:
    """
    login=test&ip=127.0.0.1
    """
    Then The error must be "nil"
    And The result must be "true"
