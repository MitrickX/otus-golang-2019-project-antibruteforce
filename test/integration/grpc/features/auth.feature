Feature: Auth
  As API client of anti brute force service
  I want check is auth allowed

  Scenario: Test auth when ip conform white list, because ip in white list
    Given "white" list with ip="127.0.0.1"

    # N auth tries, must be ok. N is limit for ip bucket
    When I call "loginLimit" times method "Auth" with params:
    """
    login=test&ip=127.0.0.1
    """
    Then The error must be "nil"
    And The result must be "true"

    # N + 1 auth try, must be ok cause of white list
    When I call method "Auth" with params:
    """
    login=test&ip=127.0.0.1
    """
    Then The error must be "nil"
    And The result must be "true"
