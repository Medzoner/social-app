describe('Feed Page', () => {
  it('should display feed posts', () => {
    cy.visit('/feed');
    cy.get('[data-cy=post], .post').should('exist');
  });
});
