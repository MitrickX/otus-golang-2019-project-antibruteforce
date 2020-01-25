Feature: Delete IP from black or white list
  As API client of anti brute force service
  I want to delete IP from black or white list

  Scenario: Delete IP from black list when ip exists
    Given "black" list with ip="127.0.0.0/24"
    When I call method "DeleteFromBlackList" with params:
    """
    ip=127.0.0.0/24
    """
    Then The error must be "nil"

  Scenario: Delete IP from black list when ip doesn't exist
    When I call method "DeleteFromBlackList" with params:
    """
    ip=127.0.0.0/24
    """
    Then The error must be "nil"

  Scenario: Delete IP from white list when ip exists
    Given "white" list with ip="127.0.0.0/24"
    When I call method "DeleteFromWhiteList" with params:
    """
    ip=127.0.0.0/24
    """
    Then The error must be "nil"

  Scenario: Delete IP from white list when ip doesnt' exist
    When I call method "DeleteFromWhiteList" with params:
    """
    ip=127.0.0.0/24
    """
    Then The error must be "nil"