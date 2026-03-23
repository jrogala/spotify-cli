Feature: Library

  Scenario: Get top tracks
    Given the Spotify API returns 5 top tracks
    When I request my top tracks
    Then the command should succeed
    And I should see 5 top tracks

  Scenario: Get top artists
    Given the Spotify API returns 3 top artists
    When I request my top artists
    Then the command should succeed
    And I should see 3 top artists

  Scenario: Get liked tracks
    Given the Spotify API returns 10 liked tracks with total 150
    When I request my liked tracks
    Then the command should succeed
    And I should see 10 liked tracks
    And the total should be 150

  Scenario: Get recently played tracks
    Given the Spotify API returns 5 recently played tracks
    When I request my recently played tracks
    Then the command should succeed
    And I should see 5 recent tracks
