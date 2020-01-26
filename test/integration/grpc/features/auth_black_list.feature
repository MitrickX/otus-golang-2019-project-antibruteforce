Feature: Auth
  As API client of anti brute force service
  I want check is auth allowed

  Scenario: Test auth when ip conform black list, because ip in black list
    Given "black" list with ip="127.0.0.1"

    # auth try, must be false, cause of black list
    When I call method "Auth" with params:
    """
    login=test&ip=127.0.0.1
    """
    Then The error must be "nil"
    And The result must be "false"

  Scenario: Test auth when ip conform black list, because there is subnet ip in black list that conform this ip
    Given "black" list with ip="127.0.0.0/24"

    # auth try, must be false, cause of black list
    When I call method "Auth" with params:
    """
    login=test&ip=127.0.0.1
    """
    Then The error must be "nil"
    And The result must be "false"