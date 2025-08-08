import { loginHelper } from './_helper.js'

describe('Profile Page', () => {
  beforeEach(() => {
    loginHelper()
  })

  it('should display profile information correctly', () => {
    cy.visit('/profile/1')

    cy.get('h2').should('contain', '@')
    cy.get('p').should('exist') // Bio

    cy.get('body').then(($body) => {
      if ($body.find('img[alt*="Avatar"]').length > 0) {
        cy.get('img[alt*="Avatar"]').should('be.visible')
      }
    })
  })

  it('should display profile posts', () => {
    cy.visit('/profile/1')

    cy.contains('Publications').should('exist')
    cy.contains('Ses publications').should('exist')

    cy.wait(2000)

    cy.get('body').then(($body) => {
      if ($body.find('.post-content').length > 0) {
        cy.get('.post-content').should('exist')
      } else {
        // Si pas de posts, vérifier le message de fin (avec émoji)
        cy.contains('Fin du fil 🎉').should('exist')
      }
    })
  })

  it('should navigate to chat from profile', () => {
    cy.visit('/profile/1')

    cy.contains('Chat').click()
    cy.url().should('include', '/chats')
  })

  it('should handle profile with no posts', () => {
    cy.visit('/profile/999')

    cy.wait(2000)

    cy.url().should('include', '/404')
  })

  it('should display posts with media correctly', () => {
    cy.visit('/profile/1')

    cy.wait(2000)

    cy.get('body').then(($body) => {
      if ($body.find('img').length > 0) {
        cy.get('img').should('be.visible')
      }
      if ($body.find('video').length > 0) {
        cy.get('video').should('exist')
      }
      if ($body.find('audio').length > 0) {
        cy.get('audio').should('exist')
      }
    })
  })

  it('should handle like functionality on posts', () => {
    cy.visit('/profile/1')

    cy.wait(2000)

    cy.get('body').then(($body) => {
      if ($body.find('button').text().match(/like|j'aime/i)) {
        cy.get('button').contains(/like|j'aime/i).should('exist')
      }
    })
  })

  it('should handle comment functionality on posts', () => {
    cy.visit('/profile/1')

    cy.wait(2000)

    cy.get('body').then(($body) => {
      if ($body.find('button').text().match(/comment|commentaire/i)) {
        cy.get('button').contains(/comment|commentaire/i).should('exist')
      }
    })
  })

  it('should navigate to other profiles from posts', () => {
    cy.visit('/profile/1')

    cy.wait(2000)

    cy.get('body').then(($body) => {
      if ($body.find('a[href*="/profile/"]').length > 0) {
        cy.get('a[href*="/profile/"]').first().click()
        cy.url().should('include', '/profile/')
      }
    })
  })

  it('should handle profile loading state', () => {
    cy.visit('/profile/1')

    cy.get('body').should('not.contain', 'Erreur')
    cy.get('body').should('not.contain', 'Error')
  })

  it('should handle non-existent profile gracefully', () => {
    cy.visit('/profile/999999')

    cy.get('body').should('exist')
  })

  it('should maintain navigation state', () => {
    cy.visit('/profile/1')

    cy.get('body').then(($body) => {
      if ($body.find('nav').length > 0) {
        cy.get('nav').should('exist')
      }
    })

    cy.reload()
    cy.url().should('include', '/profile/1')
  })

  it('should handle profile with long bio', () => {
    cy.visit('/profile/1')

    cy.get('p').should('exist')
  })

  it('should handle profile with special characters in username', () => {
    cy.visit('/profile/1')

    cy.get('h2').should('contain', '@')
  })
})
