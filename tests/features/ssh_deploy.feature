Feature: Deploy to ssh host
  In order to understand that Deployer will make release on remote host
  I want to get source files on remote machine

  Scenario: Deploy to ssh
    When I run Deployer
    Then Exit code will be zero
    And  I have folders on remote server
