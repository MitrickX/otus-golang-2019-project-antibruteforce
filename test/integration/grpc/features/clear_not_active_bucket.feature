Feature: Clear not active buckets
  As API client of anti brute force service
  I want check if not active buckets cleaning works correctly

  Scenario: Test not active buckets cleaning
    Given Clean "white" list
    And Clean "black" list

    When I call method "Auth" with params:
    """
    login=test&password=1234&ip=127.0.0.1
    """
    And Wait unit all bucket storages empty or 5 minutes

    Then The all bucket storages are empty