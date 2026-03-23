Feature: Search

  Scenario: Search for tracks, albums, and artists
    Given the Spotify API returns search results for "test" with 3 tracks and 2 albums and 1 artists
    When I search for "test" with type "track,album,artist"
    Then the command should succeed
    And I should see 3 track results
    And I should see 2 album results
    And I should see 1 artist results
