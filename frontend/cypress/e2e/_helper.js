const buildToken = (username = 'tester') => {
  const header = btoa(JSON.stringify({ alg: 'none', typ: 'JWT' }))
  const payload = btoa(
    JSON.stringify({
      sub: 1,
      username,
      role: 'user',
      verified: true,
      email: `${username}@example.com`,
      exp: Math.floor(Date.now() / 1000) + 3600
    })
  )
  return `${header}.${payload}.`
}

export const loginHelper = (username = 'tester', password = '12345') => {
  cy.clearCookies()
  cy.clearLocalStorage()

  // Stub login API pour tests e2e
  // const token = buildToken(username)
  // cy.intercept('POST', '/api/login', {
  //   statusCode: 200,
  //   body: {
  //     access_token: token,
  //     refresh_token: token,
  //     id_token: token
  //   }
  // }).as('loginRequest')

  cy.visit('/login')
  cy.get('input[type=email], input[name=email], input[name=username]')
    .should('not.be.disabled')
    .type(username)
  cy.get('input[type=password], input[name=password]').should('not.be.disabled').type(password)
  cy.get('form').submit()
  // cy.wait('@loginRequest')
  cy.wait(500)
  cy.url().should('include', '/feed')
}
