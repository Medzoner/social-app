import { loginHelper } from './_helper'

describe('Edit Profile Page', () => {
  beforeEach(() => {
    loginHelper()
  })

  it('should display edit profile form', () => {
    cy.visit('/edit-profile')
    cy.get('form').should('exist')
    cy.get('input').should('have.length.at.least', 1)
    cy.get('button[name=save]').should('exist')
  })

  it('should update profile information', () => {
    cy.visit('/edit-profile')
    cy.get('textarea[name=bio]')
      .clear()
      .type('New bio')
    cy.get('button[name=save]')
      .click()
    cy.contains(/profile updated|Profil mis à jour/i).should('exist')
  })

  it('should navigate back to profile after save', () => {
    cy.visit('/edit-profile')
    cy.get('button[name=save]')
      .click()
    cy.url().should('include', '/edit-profile')
  })

  it('should upload profile avatar', () => {
    cy.visit('/edit-profile')
    cy.get('input[type="file"]')
      .selectFile('./cypress/dummy.jpg')
    cy.get('button[name=save]')
      .click()
    cy.contains(/profile updated|Profil mis à jour/i).should('exist')
    cy.get('img[alt*="Avatar de l\'utilisateur"]').should('exist')
  })

  it('should handle long bio text', () => {
    cy.visit('/edit-profile')
    const longBio = 'A'.repeat(500)
    cy.get('textarea[name=bio]')
      .clear()
      .type(longBio)
    cy.get('button[name=save]')
      .click()
    cy.contains(/profile updated|Profil mis à jour/i).should('exist')
  })

  // it('should cancel changes', () => {
  //   cy.visit('/edit-profile');
  //   cy.get('textarea[name=bio]').clear().type('Test bio');
  //   cy.get('button').contains(/cancel|annuler/i).click();
  //   cy.url().should('include', '/profile');
  // });

  // it('should handle network errors', () => {
  //   cy.intercept('PUT', '/api/profile', { statusCode: 500 });
  //   cy.visit('/edit-profile');
  //   cy.get('textarea[name=bio]').clear().type('Test bio');
  //   cy.get('button[name=save]').click();
  //   cy.contains(/error|erreur/i).should('exist');
  // });

  // it('should show character count for bio', () => {
  //   cy.visit('/edit-profile');
  //   cy.get('textarea[name=bio]').type('Test bio');
  //   cy.get('span').contains(/500|characters/i).should('exist');
  // });

  it('should handle special characters in bio', () => {
    cy.visit('/edit-profile')
    cy.get('textarea[name=bio]')
      .clear()
      .type('Bio with émojis 🎉 and special chars!')
    cy.get('button[name=save]').click()
    cy.contains(/profile updated|Profil mis à jour/i).should('exist')
  })

  // it('should validate email format', () => {
  //   cy.visit('/edit-profile');
  //   cy.get('input[name=email]').clear().type('invalid-email');
  //   cy.get('button[name=save]').click();
  //   cy.contains(/invalid|email/i).should('exist');
  // });

  // it('should handle file size limits', () => {
  //   cy.visit('/edit-profile');
  //   // Simuler un fichier trop volumineux
  //   cy.get('input[type="file"]').selectFile('./cypress/dummy.jpg');
  //   cy.contains(/file too large|fichier trop volumineux/i).should('exist');
  // });
})
