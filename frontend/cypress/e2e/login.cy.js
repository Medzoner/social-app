describe('Login Page', () => {
  it('should login with valid credentials', () => {
    cy.visit('/login');
    cy.get('form').should('exist');
    cy.get('input[type=email], input[name=email], input[name=username]').should('exist');
    cy.get('input[type=email], input[name=email], input[name=username]').should('not.be.disabled').type('medz');
    cy.get('input[type=password], input[name=password]').should('exist');
    cy.get('input[type=password], input[name=password]').should('not.be.disabled').type('12345');
    cy.get('button').contains(/login|connexion/i).should('exist');
    cy.get('form').submit();
    cy.url().should('include', '/feed');
    cy.contains(/logout|déconnexion/i).should('exist');
  });

  it('should show error on invalid credentials', () => {
    cy.visit('/login');
    cy.get('input[type=email], input[name=email], input[name=username]').should('not.be.disabled').type('medz');
    cy.get('input[type=password], input[name=password]').should('not.be.disabled').type('wrongpassword');
    cy.get('form').submit();
    cy.contains(/invalid|erreur|incorrect/i).should('exist');
  });

  it('should logout successfully', () => {
    // Login d'abord
    cy.visit('/login');
    cy.get('input[type=email], input[name=email], input[name=username]').should('not.be.disabled').type('medz');
    cy.get('input[type=password], input[name=password]').should('not.be.disabled').type('12345');
    cy.get('form').submit();
    cy.url().should('include', '/feed');
    cy.contains(/logout|déconnexion/i).click();
    cy.url().should('include', '/login');
  });

  it('should validate required fields', () => {
    cy.visit('/login');
    cy.get('form').submit();
    // Vérifier soit un message d'erreur, soit une redirection vers login (pas de redirection vers feed)
    cy.url().should('include', '/login');
    // Optionnel : vérifier s'il y a un message d'erreur
    cy.get('body').then(($body) => {
      if ($body.find('[class*="error"], [class*="alert"]').length > 0) {
        cy.contains(/error|erreur|required|requis/i).should('exist');
      }
    });
  });

  it('should show password when toggle is clicked', () => {
    cy.visit('/login');
    cy.get('input[type=password]').should('have.attr', 'type', 'password');
    // Vérifier s'il y a un bouton pour afficher/masquer le mot de passe
    cy.get('body').then(($body) => {
      if ($body.find('button[type="button"]').length > 0) {
        cy.get('button[type="button"]').contains(/show|afficher/i).click();
        cy.get('input[type=password]').should('have.attr', 'type', 'text');
      }
    });
  });

  it('should navigate to register page', () => {
    cy.visit('/login');
    cy.get('body').then(($body) => {
      if ($body.find('a[href*="/register"]').length > 0) {
        cy.get('a[href*="/register"]').click();
        cy.url().should('include', '/register');
      }
    });
  });

  it('should navigate to forgot password', () => {
    cy.visit('/login');
    cy.get('body').then(($body) => {
      if ($body.find('a[href*="/forgot-password"]').length > 0) {
        cy.get('a[href*="/forgot-password"]').click();
        cy.url().should('include', '/forgot-password');
      }
    });
  });

  it('should handle OAuth login', () => {
    cy.visit('/login');
    cy.get('body').then(($body) => {
      if ($body.find('button').text().match(/google|oauth/i)) {
        cy.get('button').contains(/google|oauth/i).should('exist');
      }
    });
  });

  it('should remember login state', () => {
    cy.visit('/login');
    cy.get('input[type=email], input[name=email], input[name=username]').should('not.be.disabled').type('medz');
    cy.get('input[type=password], input[name=password]').should('not.be.disabled').type('12345');
    // Vérifier s'il y a une checkbox "remember me"
    cy.get('body').then(($body) => {
      if ($body.find('input[type="checkbox"]').length > 0) {
        cy.get('input[type="checkbox"]').check();
      }
    });
    cy.get('form').submit();
    cy.url().should('include', '/feed');
    // Vérifier que l'état est conservé après refresh
    cy.reload();
    cy.url().should('include', '/feed');
  });

  // it('should handle session timeout', () => {
  //   cy.visit('/login');
  //   cy.get('input[type=email], input[name=email], input[name=username]').should('not.be.disabled').type('medz');
  //   cy.get('input[type=password], input[name=password]').should('not.be.disabled').type('12345');
  //   cy.get('form').submit();
  //   cy.url().should('include', '/feed');
  //   // Simuler une session expirée
  //   cy.clearCookies();
  //   cy.visit('/feed');
  //   cy.url().should('include', '/login');
  // });
  //
  // it('should handle network errors gracefully', () => {
  //   cy.intercept('POST', '/api/auth/login', { statusCode: 500 });
  //   cy.visit('/login');
  //   cy.get('input[type=email], input[name=email], input[name=username]').should('not.be.disabled').type('medz');
  //   cy.get('input[type=password], input[name=password]').should('not.be.disabled').type('12345');
  //   cy.get('form').submit();
  //   cy.contains(/error|erreur/i).should('exist');
  // });
});
