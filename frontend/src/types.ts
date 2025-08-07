export interface Post {
  id: number
  content: string
  created_at: string
  updated_at?: string
  author_id: number
}

export interface Profile {
  id: number
  username: string
  bio?: string
  avatar_media: string | null
  posts: Post[]
}

export interface JwtUser {
  id: number
  sub: number
  username: string
  email: string
  role: string
  verified: boolean
}

export interface JwtToken {
  access_token: string
  refresh_token: string
  id_token?: string
  expires_in?: number
  token_type?: string
}

export interface AuthHeader {
  Authorization: string
  'Content-Type'?: string
}

export interface Notification {
  id: number
  type: string
  content?: string
  message: Message | string
  payload?: string
  timestamp: string
  read: boolean
  user_id: number
  receiver_id: number
}

export interface Media {
  id: number
  file_type: string
  file_path: string
}

export enum MediaType {
  Image = 'image',
  Video = 'video',
  Audio = 'audio'
}

export const MediaTypes = Object.values(MediaType)

export interface Message {
  id: number
  created_at: string
  content: string
  user_id: number
  sender_id: number
  receiver_id: number
  timestamp: Date
  read: boolean
  media?: Media[]
  sender_username?: string
  receiver_username?: string
  error: boolean
}

export interface MessageGoup {
  id: number
  messages: Message[]
  date: string
}
