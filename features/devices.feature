Feature: Device management

  Scenario: List devices
    Given the Spotify API returns devices:
      | id   | name       | type     | is_active | volume_percent |
      | dev1 | Living Room| Speaker  | true      | 50             |
      | dev2 | Kitchen    | Speaker  | false     | 30             |
    When I request the device list
    Then the command should succeed
    And I should see 2 devices
    And device "Living Room" should be active
