Feature: Add IP in black or white list
  As API client of anti brute force service
  I want to add IP in black or white list

  Scenario: Add IP in black list
    When I call method "AddInBlackList" with params:
    """
    ip=127.0.0.0/24
    """
    Then The error must be "nil"