export const loginHelper = () => {
  cy.visit('/login')
  cy.get('input[type=email], input[name=email], input[name=username]').should('not.be.disabled').type('medz')
  cy.get('input[type=password], input[name=password]').should('not.be.disabled').type('12345')
  cy.get('form').submit()
  cy.url().should('include', '/feed')
}