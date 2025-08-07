export const loginHelper = (username = 'tester', password = '12345') => {
  cy.clearCookies()
  cy.clearLocalStorage()

  cy.visit('/login')
  cy.get('input[type=email], input[name=email], input[name=username]')
    .should('not.be.disabled')
    .type(username)
  cy.get('input[type=password], input[name=password]').should('not.be.disabled').type(password)
  cy.get('form').submit()
  cy.wait(500)
  cy.url().should('include', '/feed')
}
