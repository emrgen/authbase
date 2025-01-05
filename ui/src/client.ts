import {Config, AuthbaseClient} from '@emrgen/authbase-client-ts'

const token = () => {
  return localStorage.getItem('token') || ''
}

const config = new Config(token, 'http://localhost:4001')

export const api = new AuthbaseClient(config)