Feature: Investment

  Scenario: Users can register an investment
    Given I am on the account page
    When I click the ADD button
    Then a form opens for me to describe the investment

  Scenario: User can edit a existent investment
    Given I am on the account page
    When select a existent investment on the list
    Then I can change the item values
    