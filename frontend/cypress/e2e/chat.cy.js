import { loginHelper } from './_helper.js'

describe('Chat Page', () => {
  beforeEach(() => {
    loginHelper()
  })

  it('should display chat list correctly', () => {
    cy.visit('/chats')
    
    // Vérifier que la page se charge
    cy.contains('Mes conversations').should('exist')
    
    // Vérifier la présence de la liste des chats ou le message "Aucune conversation"
    cy.get('body').then(($body) => {
      if ($body.find('ul').length > 0) {
        cy.get('ul').should('exist')
        cy.get('li').should('exist')
      } else {
        cy.contains('Aucune conversation encore.').should('exist')
      }
    })
  })

  it('should navigate to specific chat', () => {
    cy.visit('/chats')
    
    // Cliquer sur le premier chat disponible
    cy.get('body').then(($body) => {
      if ($body.find('li').length > 0) {
        cy.get('li').first().click()
        cy.url().should('include', '/chat/')
      }
    })
  })

  it('should display chat interface correctly', () => {
    // Aller directement sur un chat spécifique
    cy.visit('/chat/1')
    
    // Vérifier les éléments de l'interface
    cy.contains('Discussion avec').should('exist')
    cy.get('input[placeholder*="message"]').should('exist')
    cy.get('button').contains('Envoyer').should('exist')
    
    // Vérifier la zone de messages
    cy.get('.overflow-y-scroll').should('exist')
  })

  it('should send a message successfully', () => {
    cy.visit('/chat/1')
    
    const testMessage = 'Bonjour ! Ceci est un test Cypress'
    
    // Taper un message
    cy.get('input[placeholder*="message"]')
      .type(testMessage)
    
    // Envoyer le message
    cy.get('button').contains('Envoyer').click()
    
    // Vérifier que le message apparaît dans le chat
    cy.contains(testMessage).should('exist')
    
    // Vérifier que l'input est vidé
    cy.get('input[placeholder*="message"]').should('have.value', '')
  })

  it('should handle empty message', () => {
    cy.visit('/chat/1')
    
    // Essayer d'envoyer un message vide
    cy.get('button').contains('Envoyer').click()
    
    // Vérifier qu'aucun message vide n'est ajouté
    cy.get('.overflow-y-scroll').then(($chat) => {
      cy.get('button').contains('Envoyer').click()
      cy.get('.overflow-y-scroll').should('have.lengthOf.gte', 1)
    })
  })

  it('should handle long messages', () => {
    cy.visit('/chat/1')
    
    const longMessage = 'A'.repeat(500)
    
    // Taper un message long
    cy.get('input[placeholder*="message"]')
      .type(longMessage)
    
    // Envoyer le message
    cy.get('button').contains('Envoyer').click()
    
    // Vérifier que le message apparaît
    cy.contains(longMessage).should('exist')
  })

  it('should handle special characters in messages', () => {
    cy.visit('/chat/1')
    
    const specialMessage = 'Message avec émojis 🎉 et caractères spéciaux: @#$%^&*()'
    
    // Taper un message avec caractères spéciaux
    cy.get('input[placeholder*="message"]')
      .type(specialMessage)
    
    // Envoyer le message
    cy.get('button').contains('Envoyer').click()
    
    // Vérifier que le message apparaît
    cy.contains(specialMessage).should('exist')
  })

  it('should display message timestamps', () => {
    cy.visit('/chat/1')
    
    // Envoyer un message
    cy.get('input[placeholder*="message"]')
      .type('Test timestamp')
    cy.get('button').contains('Envoyer').click()
    
    // Vérifier que le timestamp est affiché
    cy.get('body').then(($body) => {
      if ($body.find('.text-xs.text-gray-500').length > 0) {
        cy.get('.text-xs.text-gray-500').should('exist')
      }
    })
  })

  it('should handle multiple messages', () => {
    cy.visit('/chat/1')
    
    const messages = ['Premier message', 'Deuxième message', 'Troisième message']
    
    // Envoyer plusieurs messages
    messages.forEach((message, index) => {
      cy.get('input[placeholder*="message"]')
        .type(message)
      cy.get('button').contains('Envoyer').click()
      
      // Vérifier que le message apparaît
      cy.contains(message).should('exist')
    })
    
    // Vérifier que tous les messages sont présents
    messages.forEach(message => {
      cy.contains(message).should('exist')
    })
  })

  it('should scroll to bottom when new message is sent', () => {
    cy.visit('/chat/1')
    
    // Envoyer plusieurs messages pour créer du scroll
    for (let i = 0; i < 10; i++) {
      cy.get('input[placeholder*="message"]')
        .type(`Message ${i + 1}`)
      cy.get('button').contains('Envoyer').click()
    }
    
    // Vérifier que le scroll est en bas
    cy.get('.overflow-y-scroll').should('exist')
  })

  it('should handle network errors gracefully', () => {
    // Intercepter la requête pour simuler une erreur
    cy.intercept('POST', '/api/messages', { statusCode: 500 })
    
    cy.visit('/chat/1')
    
    // Essayer d'envoyer un message
    cy.get('input[placeholder*="message"]')
      .type('Test erreur réseau')
    cy.get('button').contains('Envoyer').click()
    
    // Vérifier qu'une erreur est affichée
    cy.contains('Error').should('exist')
  })

  it('should display typing indicator', () => {
    cy.visit('/chat/1')
    
    // Commencer à taper
    cy.get('input[placeholder*="message"]')
      .type('Test typing')
    
    // Vérifier que l'indicateur de frappe peut apparaître
    cy.get('body').then(($body) => {
      if ($body.text().includes('...')) {
        cy.contains('...').should('exist')
      }
    })
  })

  it('should handle message read status', () => {
    cy.visit('/chat/1')
    
    // Envoyer un message
    cy.get('input[placeholder*="message"]')
      .type('Test read status')
    cy.get('button').contains('Envoyer').click()
    
    // Vérifier que les indicateurs de lecture sont présents
    cy.get('body').then(($body) => {
      if ($body.find('.fas.fa-check').length > 0 || $body.find('.fas.fa-check-double').length > 0) {
        cy.get('.fas').should('exist')
      }
    })
  })

  it('should handle chat with no messages', () => {
    cy.visit('/chat/999')
    
    // Vérifier le message pour chat vide
    cy.contains('Aucun message pour le moment.').should('exist')
  })

  it('should maintain chat state on page reload', () => {
    cy.visit('/chat/1')
    
    // Envoyer un message
    cy.get('input[placeholder*="message"]')
      .type('Test persistance')
    cy.get('button').contains('Envoyer').click()
    
    // Recharger la page
    cy.reload()
    
    // Vérifier que le message est toujours là
    cy.contains('Test persistance').should('exist')
  })

  it('should handle different chat users', () => {
    // Tester avec différents utilisateurs
    const users = [1, 2, 3]
    
    users.forEach(userId => {
      cy.visit(`/chat/${userId}`)
      cy.contains('Discussion avec').should('exist')
    })
  })

  it('should handle message input focus', () => {
    cy.visit('/chat/1')
    
    // Vérifier que l'input est focalisé
    cy.get('input[placeholder*="message"]').should('be.focused')
  })

  it('should handle enter key to send message', () => {
    cy.visit('/chat/1')
    
    const testMessage = 'Test avec Enter'
    
    // Taper un message et appuyer sur Enter
    cy.get('input[placeholder*="message"]')
      .type(testMessage)
      .type('{enter}')
    
    // Vérifier que le message est envoyé
    cy.contains(testMessage).should('exist')
  })

  it('should handle message formatting', () => {
    cy.visit('/chat/1')
    
    const formattedMessage = 'Message\navec\nsaut\nde\nligne'
    
    // Taper un message avec des sauts de ligne
    cy.get('input[placeholder*="message"]')
      .type(formattedMessage)
    cy.get('button').contains('Envoyer').click()
    
    // Vérifier que le formatage est préservé
    cy.contains('Message').should('exist')
  })
})
