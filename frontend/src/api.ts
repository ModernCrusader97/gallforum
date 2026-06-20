import axios from 'axios'

const api = axios.create({ baseURL: '/api' })

api.interceptors.request.use(cfg => {
  const token = localStorage.getItem('token')
  if (token) cfg.headers.Authorization = `Bearer ${token}`
  return cfg
})

export interface User { id: number; username: string; email: string }
export interface Channel { id: number; slug: string; name: string; description: string }
export interface Post {
  id: number; channel_id: number; channel_slug: string
  user_id: number | null; username: string; guest_name: string
  title: string; content: string; image_urls: string[]
  likes: number; dislikes: number; comment_count: number
  created_at: string; updated_at: string
}
export interface Comment {
  id: number; post_id: number; parent_id: number | null
  user_id: number | null; username: string; guest_name: string
  content: string; replies: Comment[]; created_at: string
}

export const auth = {
  register: (data: { username: string; email: string; password: string }) =>
    api.post('/auth/register', data),
  login: (data: { email: string; password: string }) =>
    api.post<{ token: string; user: User }>('/auth/login', data),
  me: () => api.get<User>('/auth/me'),
}

export const channels = {
  list: () => api.get<Channel[]>('/channels'),
  get: (slug: string) => api.get<Channel>(`/channels/${slug}`),
  create: (data: { slug: string; name: string; description: string }) =>
    api.post('/channels', data),
}

export const posts = {
  list: (slug: string, page = 1) =>
    api.get<{ posts: Post[]; total: number; page: number }>(`/channels/${slug}/posts?page=${page}`),
  get: (id: number) => api.get<Post>(`/posts/${id}`),
  create: (slug: string, data: object) => api.post(`/channels/${slug}/posts`, data),
  vote: (id: number, vote: 1 | -1) => api.post(`/posts/${id}/vote`, { vote }),
}

export const comments = {
  list: (postId: number) => api.get<Comment[]>(`/posts/${postId}/comments`),
  create: (postId: number, data: object) => api.post(`/posts/${postId}/comments`, data),
}

export const upload = {
  image: async (file: File) => {
    const fd = new FormData()
    fd.append('image', file)
    const res = await api.post<{ url: string }>('/upload', fd)
    return res.data.url
  }
}

export default api
