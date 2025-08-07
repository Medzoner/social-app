describe('Chat Page', () => {
  it('should display chat area', () => {
    cy.visit('/chat');
    cy.get('[data-cy=chat], .chat').should('exist');
  });
});
