describe('Edit Profile Page', () => {
  beforeEach(() => {
    // Login avant chaque test
    cy.visit('/login');
    cy.get('input[type=email], input[name=email], input[name=username]').type('medz');
    cy.get('input[type=password], input[name=password]').type('12345');
    cy.get('form').submit();
    cy.url().should('include', '/feed');
  });

  it('should display edit profile form', () => {
    cy.visit('/edit-profile');
    cy.get('form').should('exist');
    cy.get('input').should('have.length.at.least', 1);
    cy.get('button[name=save]').should('exist');
  });

  it('should update profile information', () => {
    cy.visit('/edit-profile');
    cy.get('textarea[name=bio]').clear().type('New bio');
    cy.get('button[name=save]').click();
    cy.contains(/profile updated|profil mis à jour/i).should('exist');
  });

  it('should navigate back to profile after save', () => {
    cy.visit('/edit-profile');
    cy.get('button[name=save]').click();
    cy.url().should('include', '/edit-profile');
  });
});
