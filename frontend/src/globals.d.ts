export {}

declare global {
  interface Window {
    ws: WebSocket | null
  }
}
