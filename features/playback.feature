Feature: Playback control

  Scenario: Get now-playing state
    Given the Spotify API returns a playback state with track "Bohemian Rhapsody" by "Queen" on album "A Night at the Opera"
    When I request the now-playing state
    Then the command should succeed
    And I should see track "Bohemian Rhapsody"
    And I should see artist "Queen"

  Scenario: Resume playback
    Given the Spotify API accepts play commands
    When I send a play command
    Then the command should succeed

  Scenario: Pause playback
    Given the Spotify API accepts pause commands
    When I send a pause command
    Then the command should succeed

  Scenario: Skip to next track
    Given the Spotify API accepts next commands
    When I send a next command
    Then the command should succeed

  Scenario: Skip to previous track
    Given the Spotify API accepts previous commands
    When I send a previous command
    Then the command should succeed
