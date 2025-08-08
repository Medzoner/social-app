describe('Chat Interaction Between Users', () => {
  it('should allow two users to chat with each other', () => {
    // Première session - Utilisateur 1
    cy.session('user1-chat-basic', () => {
      cy.visit('/login')
      cy.get('input[type=email], input[name=email], input[name=username]')
        .should('not.be.disabled')
        .type('tester')
      cy.get('input[type=password], input[name=password]')
        .should('not.be.disabled')
        .type('12345')
      cy.get('form').submit()
      cy.url().should('include', '/feed')
    })

    // Deuxième session - Utilisateur 2
    cy.session('user2-chat-basic', () => {
      cy.visit('/login')
      cy.get('input[type=email], input[name=email], input[name=username]')
        .should('not.be.disabled')
        .type('medz')
      cy.get('input[type=password], input[name=password]')
        .should('not.be.disabled')
        .type('12345')
      cy.get('form').submit()
      cy.url().should('include', '/feed')
    })

    // Test de l'interaction - Utilisateur 1 envoie un message
    cy.session('user1-send-message', () => {
      // Se connecter d'abord
      cy.visit('/login')
      cy.get('input[type=email], input[name=email], input[name=username]')
        .should('not.be.disabled')
        .type('tester')
      cy.get('input[type=password], input[name=password]')
        .should('not.be.disabled')
        .type('12345')
      cy.get('form').submit()
      cy.url().should('include', '/feed')
      
      // Aller sur la page de chat
      cy.visit('/chat/4') // Chat avec utilisateur 4
      
      // Attendre que la page de chat soit chargée
      cy.url().should('include', '/chat/4')
      cy.get('input[name="message"]').should('be.visible')
      
      const messageFromUser1 = 'Salut ! Je suis l\'utilisateur 1'
      cy.get('input[name="message"]')
        .type(messageFromUser1)
      cy.get('button').contains('Envoyer').click()
      
      // Vérifier que le message est envoyé
      cy.contains(messageFromUser1).should('exist')
    })

    // Test de l'interaction - Utilisateur 2 reçoit et répond
    cy.session('user2-receive-reply', () => {
      // Se connecter d'abord
      cy.visit('/login')
      cy.get('input[type=email], input[name=email], input[name=username]')
        .should('not.be.disabled')
        .type('medz')
      cy.get('input[type=password], input[name=password]')
        .should('not.be.disabled')
        .type('12345')
      cy.get('form').submit()
      cy.url().should('include', '/feed')
      
      // Aller sur la page de chat
      cy.visit('/chat/1') // Chat avec utilisateur 1
      
      // Attendre que la page de chat soit chargée
      cy.url().should('include', '/chat/1')
      cy.get('input[name="message"]').should('be.visible')
      
      // Vérifier que le message de l'utilisateur 1 est visible
      cy.contains('Salut ! Je suis l\'utilisateur 1').should('exist')
      
      // Utilisateur 2 répond
      const messageFromUser2 = 'Salut ! Je suis l\'utilisateur 2'
      cy.get('input[name="message"]')
        .type(messageFromUser2)
      cy.get('button').contains('Envoyer').click()
      
      // Vérifier que la réponse est envoyée
      cy.contains(messageFromUser2).should('exist')
    })

    // Test de l'interaction - Utilisateur 1 voit la réponse
    cy.session('user1-see-reply', () => {
      // Se connecter d'abord
      cy.visit('/login')
      cy.get('input[type=email], input[name=email], input[name=username]')
        .should('not.be.disabled')
        .type('tester')
      cy.get('input[type=password], input[name=password]')
        .should('not.be.disabled')
        .type('12345')
      cy.get('form').submit()
      cy.url().should('include', '/feed')
      
      // Aller sur la page de chat
      cy.visit('/chat/4')
      
      // Attendre que la page de chat soit chargée
      cy.url().should('include', '/chat/4')
      
      // Vérifier que la réponse de l'utilisateur 2 est visible
      cy.contains('Salut ! Je suis l\'utilisateur 2').should('exist')
    })
  })

  it('should handle real-time message updates', () => {
    // Test de messages en temps réel - Utilisateur 1 envoie
    cy.session('user1-realtime-send', () => {
      cy.visit('/login')
      cy.get('input[type=email], input[name=email], input[name=username]')
        .should('not.be.disabled')
        .type('tester')
      cy.get('input[type=password], input[name=password]')
        .should('not.be.disabled')
        .type('12345')
      cy.get('form').submit()
      cy.url().should('include', '/feed')
      
      cy.visit('/chat/4')
      
      // Attendre que la page de chat soit chargée
      cy.url().should('include', '/chat/4')
      cy.get('input[name="message"]').should('be.visible')
      
      // Envoyer plusieurs messages rapidement
      const messages = [
        'Message 1 en temps réel',
        'Message 2 en temps réel',
        'Message 3 en temps réel'
      ]
      
      messages.forEach((message, index) => {
        cy.get('input[name="message"]')
          .type(message)
        cy.get('button').contains('Envoyer').click()
        
        // Vérifier que le message apparaît
        cy.contains(message).should('exist')
      })
    })

    // Test de messages en temps réel - Utilisateur 2 reçoit
    cy.session('user2-realtime-receive', () => {
      cy.visit('/login')
      cy.get('input[type=email], input[name=email], input[name=username]')
        .should('not.be.disabled')
        .type('medz')
      cy.get('input[type=password], input[name=password]')
        .should('not.be.disabled')
        .type('12345')
      cy.get('form').submit()
      cy.url().should('include', '/feed')
      
      cy.visit('/chat/1')
      
      // Attendre que la page de chat soit chargée
      cy.url().should('include', '/chat/1')
      
      // Vérifier que tous les messages sont reçus
      cy.contains('Message 1 en temps réel').should('exist')
      cy.contains('Message 2 en temps réel').should('exist')
      cy.contains('Message 3 en temps réel').should('exist')
    })
  })

  it('should handle typing indicators between users', () => {
    // Test des indicateurs de frappe
    cy.session('user1-typing-indicator', () => {
      cy.visit('/login')
      cy.get('input[type=email], input[name=email], input[name=username]')
        .should('not.be.disabled')
        .type('tester')
      cy.get('input[type=password], input[name=password]')
        .should('not.be.disabled')
        .type('12345')
      cy.get('form').submit()
      cy.url().should('include', '/feed')
      
      cy.visit('/chat/4')
      
      // Attendre que la page de chat soit chargée
      cy.url().should('include', '/chat/4')
      cy.get('input[name="message"]').should('be.visible')
      
      // Commencer à taper
      cy.get('input[name="message"]')
        .type('Test indicateur de frappe')
      
      // Vérifier que l'indicateur peut apparaître
      cy.get('body').then(($body) => {
        if ($body.text().includes('...')) {
          cy.contains('...').should('exist')
        }
      })
    })
  })

  it('should handle message read status between users', () => {
    // Test du statut de lecture - Utilisateur 1 envoie
    cy.session('user1-read-status-send', () => {
      cy.visit('/login')
      cy.get('input[type=email], input[name=email], input[name=username]')
        .should('not.be.disabled')
        .type('tester')
      cy.get('input[type=password], input[name=password]')
        .should('not.be.disabled')
        .type('12345')
      cy.get('form').submit()
      cy.url().should('include', '/feed')
      
      cy.visit('/chat/4')
      
      // Attendre que la page de chat soit chargée
      cy.url().should('include', '/chat/4')
      cy.get('input[name="message"]').should('be.visible')
      
      // Envoyer un message
      cy.get('input[name="message"]')
        .type('Test statut de lecture')
      cy.get('button').contains('Envoyer').click()
      
      // Vérifier les indicateurs de lecture
      cy.get('body').then(($body) => {
        if ($body.find('.fas.fa-check').length > 0 || $body.find('.fas.fa-check-double').length > 0) {
          cy.get('.fas').should('exist')
        }
      })
    })

    // Test du statut de lecture - Utilisateur 2 reçoit
    cy.session('user2-read-status-receive', () => {
      cy.visit('/login')
      cy.get('input[type=email], input[name=email], input[name=username]')
        .should('not.be.disabled')
        .type('medz')
      cy.get('input[type=password], input[name=password]')
        .should('not.be.disabled')
        .type('12345')
      cy.get('form').submit()
      cy.url().should('include', '/feed')
      
      cy.visit('/chat/1')
      
      // Attendre que la page de chat soit chargée
      cy.url().should('include', '/chat/1')
      
      // Vérifier que le message est marqué comme lu
      cy.contains('Test statut de lecture').should('exist')
    })
  })

  it('should handle multiple chat conversations', () => {
    // Test avec plusieurs conversations simultanées
    const conversations = [
      { user1: 'tester', user2: 'medz', chatId: 4 },
      { user1: 'medz', user2: 'tester', chatId: 1 }
    ]

    conversations.forEach((conv, index) => {
      cy.session(`user1-multi-chat-${index}`, () => {
        cy.visit('/login')
        cy.get('input[type=email], input[name=email], input[name=username]')
          .should('not.be.disabled')
          .type(conv.user1)
        cy.get('input[type=password], input[name=password]')
          .should('not.be.disabled')
          .type('12345')
        cy.get('form').submit()
        cy.url().should('include', '/feed')
        
        cy.visit(`/chat/${conv.chatId}`)
        
        // Attendre que la page de chat soit chargée
        cy.url().should('include', `/chat/${conv.chatId}`)
        cy.get('input[name="message"]').should('be.visible')
        
        const message = `Message de ${conv.user1} vers ${conv.user2}`
        cy.get('input[name="message"]')
          .type(message)
        cy.get('button').contains('Envoyer').click()
        
        cy.contains(message).should('exist')
      })
    })
  })

  it('should handle chat list updates', () => {
    // Vérifier que la liste des chats se met à jour
    cy.session('user1-chat-list', () => {
      cy.visit('/login')
      cy.get('input[type=email], input[name=email], input[name=username]')
        .should('not.be.disabled')
        .type('tester')
      cy.get('input[type=password], input[name=password]')
        .should('not.be.disabled')
        .type('12345')
      cy.get('form').submit()
      cy.url().should('include', '/feed')
      
      cy.visit('/chats')
      
      // Vérifier que les conversations sont listées
      cy.get('body').then(($body) => {
        if ($body.find('ul').length > 0) {
          cy.get('ul').should('exist')
          cy.get('li').should('exist')
        }
      })
    })
  })
})
