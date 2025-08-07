import { loginHelper } from './_helper.js'

describe('Feed Page', () => {
  beforeEach(() => {
    loginHelper()
  });

  it('should display feed posts', () => {
    cy.visit('/feed');
    cy.get('li').should('exist');
  });

  it('should handle modal close', () => {
    cy.visit('/feed');
    cy.get('button').contains(/créer un post/i).click();
    cy.get('textarea').should('be.visible');
    cy.get('button[name="close-button"]').click();
    cy.get('.fixed.inset-0').should('not.exist');
    cy.wait(500)
    cy.get('#post-modal').should('not.exist');
  });

  it('should create a new post', () => {
    cy.visit('/feed');
    cy.get('button').contains(/créer un post/i).click();
    cy.get('textarea').should('be.visible').first().type('Test post from Cypress');
    cy.get('button').contains(/publier/i).click();
    cy.contains('Test post from Cypress').should('exist');
  });

  it('should create a post with image upload', () => {
    cy.visit('/feed');
    cy.get('button').contains(/créer un post/i).click();
    cy.get('textarea').should('be.visible').first().type('Test post with image');
    
    // Upload d'une image
    cy.get('input[type="file"]').selectFile('./cypress/dummy.jpg');
    cy.get('button').contains(/publier/i).click();
    cy.contains('Test post with image').should('exist');
  });

  it('should search posts', () => {
    cy.visit('/feed');
    cy.get('input[placeholder*="search"]').type('test');
    cy.get('li').should('contain', 'test');
  });

  it('should like a post', () => {
    cy.visit('/feed');
    cy.get('button').contains(/j'aime/i).first().click();
    cy.get('button').contains(/je n'aime plus/i).should('exist');
  });

  it('should unlike a post', () => {
    cy.visit('/feed');
    // Like d'abord
    cy.get('button').contains(/j'aime/i).first().click();
    // Puis unlike
    cy.get('button').contains(/je n'aime plus/i).first().click();
    cy.get('button').contains(/j'aime/i).should('exist');
  });

  it('should comment on a post', () => {
    cy.visit('/feed');
    cy.get('input[placeholder*="commentaire"]').first().type('Test comment');
    cy.get('button').contains(/envoyer/i).first().click();
    cy.contains('Test comment').should('exist');
  });

  it('should show comment count', () => {
    cy.visit('/feed');
    cy.get('button').contains(/afficher les commentaires/i).should('exist');
  });

  it('should load more comments', () => {
    cy.visit('/feed');
    cy.get('button').contains(/afficher les commentaires/i).first().click();
    cy.get('ul').should('be.visible');
  });

  it('should display user avatars', () => {
    cy.visit('/feed');
    cy.get('img[alt*="Avatar"]').should('exist');
  });

  it('should navigate to user profile from post', () => {
    cy.visit('/feed');
    cy.get('a[href*="/profile/"]').first().click();
    cy.url().should('include', '/profile/');
  });

  it('should handle pagination', () => {
    cy.visit('/feed');
    // Scroll pour déclencher le chargement de plus de posts
    cy.scrollTo('bottom');
    cy.get('li').should('have.length.greaterThan', 1);
  });

  it('should show loading state', () => {
    cy.visit('/feed');
    cy.get('button').contains(/créer un post/i).click();
    cy.get('textarea').should('be.visible').first().type('Loading test post');
    cy.get('button').contains(/publier/i).click();
    cy.contains('Post publié').should('exist');
  });
});
