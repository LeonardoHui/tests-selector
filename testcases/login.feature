Feature: Login

  Scenario: Login page is accessible from homepage
    Given I am on the homepage
    When I click the LOG_IN button
    Then I am redirected to the login page

  Scenario: Registered users can login
    Given I am on the login pagin   
    When I enter my email and password
    Then I am redirected to my account dashboard

  Scenario: Users cannot log in with invalid credentions
    Given I am on the login pagin   
    When I enter an invalid email and password
    Then I get an error message
