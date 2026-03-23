Feature: Volume control

  Scenario: Get current volume
    Given the Spotify API returns volume 65 on device "Living Room"
    When I request the current volume
    Then the command should succeed
    And I should see volume 65

  Scenario: Set volume
    Given the Spotify API accepts volume commands
    When I set the volume to 50
    Then the command should succeed
