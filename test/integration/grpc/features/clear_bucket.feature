Feature: Clear bucket
  As API client of anti brute force service
  I want to clear bucket

  Scenario: Clear bucket by login when bucket exists
    Given bucket for
    """
    login=test
    """
    When I call method "ClearBucket" with params:
    """
    login=test
    """
    Then The error must be "nil"

  Scenario: Clear bucket by login when bucket doesn't exist
    When I call method "ClearBucket" with params:
    """
    login=test2
    """
    Then The error must be "nil"

  Scenario: Clear bucket by password when bucket exists
    Given bucket for
    """
    password=1234
    """
    When I call method "ClearBucket" with params:
    """
    password=1234
    """
    Then The error must be "nil"

  Scenario: Clear bucket by password when bucket doesn't exist
    When I call method "ClearBucket" with params:
    """
    password=5678
    """
    Then The error must be "nil"

  Scenario: Clear bucket by ip when bucket exists
    Given bucket for
    """
    ip=127.0.0.1
    """
    When I call method "ClearBucket" with params:
    """
    ip=127.0.0.1
    """
    Then The error must be "nil"

  Scenario: Clear bucket by ip when bucket doesn't exist
    When I call method "ClearBucket" with params:
    """
    ip=127.0.0.2
    """
    Then The error must be "nil"