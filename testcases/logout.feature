Feature: Login

  Scenario: Logout page is accessible from account page
    Given I am on the account page
    When I click the LOG_OUT button
    Then I am redirected to the homepage

  Scenario: Logged out users cannot access the account page
    Given I logged out  
    When I type the ACCOUNT_PAGE address
    Then I get an error message
    