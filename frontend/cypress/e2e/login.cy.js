describe('Login Page', () => {
  it('should login with valid credentials', () => {
    cy.visit('/login');
    cy.get('form').should('exist');
    cy.get('input[type=email], input[name=email], input[name=username]').should('exist');
    cy.get('input[type=email], input[name=email], input[name=username]').type('medz');
    cy.get('input[type=password], input[name=password]').should('exist');
    cy.get('input[type=password], input[name=password]').type('12345');
    cy.get('button').contains(/login|connexion/i).should('exist');
    cy.get('form').submit();
    cy.url().should('include', '/feed');
    cy.contains(/logout|déconnexion/i).should('exist');
  });

  it('should show error on invalid credentials', () => {
    cy.visit('/login');
    cy.get('input[type=email], input[name=email], input[name=username]').type('medz');
    cy.get('input[type=password], input[name=password]').type('wrongpassword');
    cy.get('form').submit();
    cy.contains(/invalid|erreur|incorrect/i).should('exist');
  });

  it('should logout successfully', () => {
    // Login d'abord
    cy.visit('/login');
    cy.get('input[type=email], input[name=email], input[name=username]').type('medz');
    cy.get('input[type=password], input[name=password]').type('12345');
    cy.get('form').submit();
    cy.url().should('include', '/feed');
    cy.contains(/logout|déconnexion/i).click();
    cy.url().should('include', '/login');
  });
});
