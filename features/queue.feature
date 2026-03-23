Feature: Queue management

  Scenario: List queue
    Given the Spotify API returns a queue with current track "Current Song" and 3 queued tracks
    When I request the queue
    Then the command should succeed
    And the now-playing track should be "Current Song"
    And the queue should have 3 tracks

  Scenario: Add to queue
    Given the Spotify API accepts add-to-queue commands
    When I add "spotify:track:abc123" to the queue
    Then the command should succeed
