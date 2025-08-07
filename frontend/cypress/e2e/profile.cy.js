describe('Profile Page', () => {
  it('should display user profile', () => {
    cy.visit('/profile');
    cy.get('[data-cy=profile], .profile').should('exist');
  });
});
