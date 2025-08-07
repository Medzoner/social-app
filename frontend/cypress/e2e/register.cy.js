describe('Register Page', () => {
  it('should display register form', () => {
    cy.visit('/register');
    cy.get('form').should('exist');
    cy.get('input').should('have.length.at.least', 2);
    cy.get('button').contains(/register|inscription/i).should('exist');
  });
});
