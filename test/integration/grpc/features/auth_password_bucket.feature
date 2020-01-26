Feature: Auth
  As API client of anti brute force service
  I want check is auth allowed

  Scenario: Test auth when password bucket is overflowing
    Given Clean bucket for
    """
    password=1234&ip=127.0.0.1
    """
    And Clean "white" list
    And Clean "black" list

    # N auth tries with different random login, must be ok. N is limit for password bucket
    When I call "passwordLimit" times method "Auth" with params:
    """
    login=random&password=1234&ip=127.0.0.1
    """
    Then The error must be "nil"
    And The result must be "true"

    # N + 1 auth try with same password, must be not ok cause of overflowing
    When I call method "Auth" with params:
    """
    login=random&password=1234&ip=127.0.0.1
    """
    Then The error must be "nil"
    And The result must be "false"

    # 1 auth try for different password, must be ok cause new password - new bucket
    When I call method "Auth" with params:
    """
    login=random&password=4567&ip=127.0.0.1
    """
    Then The error must be "nil"
    And The result must be "true"

  Scenario: Test auth when password bucket is not overflowing
    Given Clean bucket for
    """
    password=1234&ip=127.0.0.1
    """
    And Clean "white" list
    And Clean "black" list

    # N auth tries with different random login but with same password must be ok. N is limit for password bucket
    When I call "passwordLimit" times method "Auth" with params:
    """
    login=random&password=1234&ip=127.0.0.1
    """
    Then The error must be "nil"
    And The result must be "true"

    # wait 1 minute
    When Wait 1 minute

    # N auth tries with different random login but with same password, must be ok. N is limit for password bucket
    And I call "passwordLimit" times method "Auth" with params:
    """
    login=random&password=1234&ip=127.0.0.1
    """
    Then The error must be "nil"
    And The result must be "true"
