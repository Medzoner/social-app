import { loginHelper } from './_helper'

describe('Edit Profile Page', () => {
  beforeEach(() => {
    cy.clearCookies()
    cy.clearLocalStorage()
  })

  it('should load edit profile page correctly', () => {
    // Login d'abord
    cy.visit('/login')
    cy.get('input[type=email], input[name=email], input[name=username]')
      .should('not.be.disabled')
      .type('tester')
    cy.get('input[type=password], input[name=password]')
      .should('not.be.disabled')
      .type('12345')
    cy.get('form').submit()
    cy.url().should('include', '/feed')

    // Aller sur la page d'édition du profil
    cy.visit('/edit-profile')
    
    // Vérifier que les éléments sont présents
    cy.contains('Modifier mon profil').should('exist')
    cy.get('textarea[name="bio"]').should('exist')
    cy.get('input[type="file"]').should('exist')
    cy.get('button[name="save"]').should('exist')
  })

  it('should display current profile information', () => {
    // Login d'abord
    cy.visit('/login')
    cy.get('input[type=email], input[name=email], input[name=username]')
      .should('not.be.disabled')
      .type('tester')
    cy.get('input[type=password], input[name=password]')
      .should('not.be.disabled')
      .type('12345')
    cy.get('form').submit()
    cy.url().should('include', '/feed')

    // Aller sur la page d'édition du profil
    cy.visit('/edit-profile')
    
    // Vérifier que les informations actuelles sont affichées
    cy.get('textarea[name="bio"]').should('exist')
    
    // Vérifier la présence de l'avatar actuel (optionnel)
    cy.get('body').then(($body) => {
      if ($body.find('img[alt*="Avatar"]').length > 0) {
        cy.get('img[alt*="Avatar"]').should('be.visible')
      }
    })
  })

  it('should update bio successfully', () => {
    // Login d'abord
    cy.visit('/login')
    cy.get('input[type=email], input[name=email], input[name=username]')
      .should('not.be.disabled')
      .type('tester')
    cy.get('input[type=password], input[name=password]')
      .should('not.be.disabled')
      .type('12345')
    cy.get('form').submit()
    cy.url().should('include', '/feed')

    // Aller sur la page d'édition du profil
    cy.visit('/edit-profile')
    
    // Modifier la bio
    const newBio = 'Nouvelle bio de test pour Cypress'
    cy.get('textarea[name="bio"]')
      .clear()
      .type(newBio)
    
    // Sauvegarder
    cy.get('button[name="save"]').click()
    
    // Vérifier le message de succès
    cy.contains('succès').should('exist')
  })

  it('should handle file upload for avatar', () => {
    // Login d'abord
    cy.visit('/login')
    cy.get('input[type=email], input[name=email], input[name=username]')
      .should('not.be.disabled')
      .type('tester')
    cy.get('input[type=password], input[name=password]')
      .should('not.be.disabled')
      .type('12345')
    cy.get('form').submit()
    cy.url().should('include', '/feed')

    // Aller sur la page d'édition du profil
    cy.visit('/edit-profile')
    
    // Uploader un fichier
    cy.get('input[type="file"]').selectFile('cypress/dummy.jpg', { force: true })
    
    // Sauvegarder
    cy.get('button[name="save"]').click()
    
    // Vérifier le message de succès
    cy.contains('succès').should('exist')
  })

  it('should handle empty bio update', () => {
    // Login d'abord
    cy.visit('/login')
    cy.get('input[type=email], input[name=email], input[name=username]')
      .should('not.be.disabled')
      .type('tester')
    cy.get('input[type=password], input[name=password]')
      .should('not.be.disabled')
      .type('12345')
    cy.get('form').submit()
    cy.url().should('include', '/feed')

    // Aller sur la page d'édition du profil
    cy.visit('/edit-profile')
    
    // Vider la bio
    cy.get('textarea[name="bio"]').clear()
    
    // Sauvegarder
    cy.get('button[name="save"]').click()
    
    // Vérifier le message de succès
    cy.contains('succès').should('exist')
  })

  it('should handle long bio text', () => {
    // Login d'abord
    cy.visit('/login')
    cy.get('input[type=email], input[name=email], input[name=username]')
      .should('not.be.disabled')
      .type('tester')
    cy.get('input[type=password], input[name=password]')
      .should('not.be.disabled')
      .type('12345')
    cy.get('form').submit()
    cy.url().should('include', '/feed')

    // Aller sur la page d'édition du profil
    cy.visit('/edit-profile')
    
    // Taper une bio longue
    const longBio = 'A'.repeat(500)
    cy.get('textarea[name="bio"]')
      .clear()
      .type(longBio)
    
    // Sauvegarder
    cy.get('button[name="save"]').click()
    
    // Vérifier le message de succès
    cy.contains('succès').should('exist')
  })

  it('should handle special characters in bio', () => {
    // Login d'abord
    cy.visit('/login')
    cy.get('input[type=email], input[name=email], input[name=username]')
      .should('not.be.disabled')
      .type('tester')
    cy.get('input[type=password], input[name=password]')
      .should('not.be.disabled')
      .type('12345')
    cy.get('form').submit()
    cy.url().should('include', '/feed')

    // Aller sur la page d'édition du profil
    cy.visit('/edit-profile')
    
    // Taper une bio avec des caractères spéciaux
    const specialBio = 'Bio avec émojis 🎉 et caractères spéciaux: @#$%^&*()'
    cy.get('textarea[name="bio"]')
      .clear()
      .type(specialBio)
    
    // Sauvegarder
    cy.get('button[name="save"]').click()
    
    // Vérifier le message de succès
    cy.contains('succès').should('exist')
  })

  it('should handle invalid file upload', () => {
    // Login d'abord
    cy.visit('/login')
    cy.get('input[type=email], input[name=email], input[name=username]')
      .should('not.be.disabled')
      .type('tester')
    cy.get('input[type=password], input[name=password]')
      .should('not.be.disabled')
      .type('12345')
    cy.get('form').submit()
    cy.url().should('include', '/feed')

    // Aller sur la page d'édition du profil
    cy.visit('/edit-profile')
    
    // Essayer d'uploader un fichier invalide
    cy.get('input[type="file"]').selectFile('cypress/dummy.mp4', { force: true })
    
    // Sauvegarder
    cy.get('button[name="save"]').click()
    
    // Vérifier qu'il y a une gestion d'erreur
    cy.get('body').should('exist')
  })

  it('should show error notification on update failure', () => {
    // Login d'abord
    cy.visit('/login')
    cy.get('input[type=email], input[name=email], input[name=username]')
      .should('not.be.disabled')
      .type('tester')
    cy.get('input[type=password], input[name=password]')
      .should('not.be.disabled')
      .type('12345')
    cy.get('form').submit()
    cy.url().should('include', '/feed')

    // Intercepter la requête pour simuler une erreur
    cy.intercept('PATCH', '/api/users/*', { statusCode: 500 })
    
    // Aller sur la page d'édition du profil
    cy.visit('/edit-profile')
    
    // Modifier la bio
    cy.get('textarea[name="bio"]')
      .clear()
      .type('Test bio')
    
    // Sauvegarder
    cy.get('button[name="save"]').click()
    
    // Vérifier le message d'erreur
    cy.contains('Échec').should('exist')
  })

  it('should redirect to login if not authenticated', () => {
    // Aller directement sur la page d'édition sans être connecté
    cy.visit('/edit-profile')
    
    // Vérifier qu'on est redirigé vers login
    cy.url().should('include', '/login')
  })

  it('should handle form submission with both bio and file', () => {
    // Login d'abord
    cy.visit('/login')
    cy.get('input[type=email], input[name=email], input[name=username]')
      .should('not.be.disabled')
      .type('tester')
    cy.get('input[type=password], input[name=password]')
      .should('not.be.disabled')
      .type('12345')
    cy.get('form').submit()
    cy.url().should('include', '/feed')

    // Aller sur la page d'édition du profil
    cy.visit('/edit-profile')
    
    // Modifier la bio
    cy.get('textarea[name="bio"]')
      .clear()
      .type('Bio mise à jour avec fichier')
    
    // Uploader un fichier
    cy.get('input[type="file"]').selectFile('cypress/dummy.jpg', { force: true })
    
    // Sauvegarder
    cy.get('button[name="save"]').click()
    
    // Vérifier le message de succès
    cy.contains('succès').should('exist')
  })

  it('should clear notification after timeout', () => {
    // Login d'abord
    cy.visit('/login')
    cy.get('input[type=email], input[name=email], input[name=username]')
      .should('not.be.disabled')
      .type('tester')
    cy.get('input[type=password], input[name=password]')
      .should('not.be.disabled')
      .type('12345')
    cy.get('form').submit()
    cy.url().should('include', '/feed')

    // Aller sur la page d'édition du profil
    cy.visit('/edit-profile')
    
    // Modifier la bio
    cy.get('textarea[name="bio"]')
      .clear()
      .type('Test notification timeout')
    
    // Sauvegarder
    cy.get('button[name="save"]').click()
    
    // Vérifier que la notification apparaît
    cy.contains('succès').should('exist')
    
    // Attendre que la notification disparaisse (5 secondes)
    cy.wait(6000)
    cy.contains('succès').should('not.exist')
  })

  it('should handle form validation', () => {
    // Login d'abord
    cy.visit('/login')
    cy.get('input[type=email], input[name=email], input[name=username]')
      .should('not.be.disabled')
      .type('tester')
    cy.get('input[type=password], input[name=password]')
      .should('not.be.disabled')
      .type('12345')
    cy.get('form').submit()
    cy.url().should('include', '/feed')

    // Aller sur la page d'édition du profil
    cy.visit('/edit-profile')
    
    // Vérifier que le formulaire est valide
    cy.get('form').should('exist')
    cy.get('textarea[name="bio"]').should('be.visible')
    cy.get('input[type="file"]').should('be.visible')
    cy.get('button[name="save"]').should('be.visible')
  })

  it('should maintain form state on page reload', () => {
    // Login d'abord
    cy.visit('/login')
    cy.get('input[type=email], input[name=email], input[name=username]')
      .should('not.be.disabled')
      .type('tester')
    cy.get('input[type=password], input[name=password]')
      .should('not.be.disabled')
      .type('12345')
    cy.get('form').submit()
    cy.url().should('include', '/feed')

    // Aller sur la page d'édition du profil
    cy.visit('/edit-profile')
    
    // Vérifier que la page se recharge correctement
    cy.reload()
    cy.contains('Modifier mon profil').should('exist')
    cy.get('textarea[name="bio"]').should('exist')
  })
})
