Feature: Run deployer with options
  In order to understand that Deployer will exit as we wait

  Scenario: Run for test
    When I run Deployer for test
    Then Exit code will be zero

  Scenario: Run without target
    When I run Deployer without target
    Then Exit code will not zero

  Scenario: Run with wrong target
    When I run Deployer with wrong target
    Then Exit code will not zero