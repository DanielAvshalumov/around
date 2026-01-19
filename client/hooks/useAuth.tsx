import { createContext, ReactNode, useState } from 'react'
import { User } from '../lib/auth'

interface AuthContextType {
    user: User | null;
    token: string | null;
    loading: boolean;
    login: (provider: 'google' | 'microsoft') => void;
    logout: () => Promise<void>;
    setToken: (token: string) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: {children: ReactNode}) {
    const [user, setUser] = useState<User | null>(null)
    const [password, setPassword] = useState<string>("")
    const [loading, setLoading] = useState(true)

    const login = (provider: 'google' | 'microsoft') => {
        const loginUrl = 'http://localhost/' + `/api/auth/${provider}`
    }
}