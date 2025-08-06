import { ChatOllama } from '@langchain/ollama'
import { HumanMessage, SystemMessage } from '@langchain/core/messages'

const model = new ChatOllama({
  modelName: 'llama3',
  temperature: 0,
  configuration: {
    baseURL: 'http://localhost:11434/v1'
  }
})

const messages = [
  new SystemMessage(
    `Your name is Helix. You are a simple helpful chatbot. You answer user queries in a friendly manner. 
     Make your responses simpler.`
  )
]

export const sendMessage = async (message) => {
  messages.push(new HumanMessage(message))
  const response = await model.invoke(messages)
  return response.content
}
