describe('Notifications Page', () => {
  it('should display notifications', () => {
    cy.visit('/notifications');
    cy.get('[data-cy=notification], .notification').should('exist');
  });
});
