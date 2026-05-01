import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { AuthProvider } from './components/AuthProvider'
import { ProtectedRoute, SuperAdminRoute } from './components/ProtectedRoute'
import Layout from './components/Layout'
import LoginPage from './pages/LoginPage'
import DashboardPage from './pages/DashboardPage'
import PatientListPage from './pages/PatientListPage'
import PatientDetailPage from './pages/PatientDetailPage'
import PatientFormPage from './pages/PatientFormPage'
import ExportPage from './pages/ExportPage'
import UserManagementPage from './pages/UserManagementPage'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      staleTime: 30_000,
    },
  },
})

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <BrowserRouter>
          <Routes>
            <Route path="/login" element={<LoginPage />} />

            <Route element={<ProtectedRoute />}>
              <Route element={<Layout />}>
                <Route index element={<DashboardPage />} />
                <Route path="patients" element={<PatientListPage />} />
                <Route path="patients/new" element={<PatientFormPage />} />
                <Route path="patients/:id" element={<PatientDetailPage />} />
                <Route path="patients/:id/edit" element={<PatientFormPage />} />
                <Route path="export" element={<ExportPage />} />

                <Route element={<SuperAdminRoute />}>
                  <Route path="users" element={<UserManagementPage />} />
                </Route>
              </Route>
            </Route>

            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </BrowserRouter>
      </AuthProvider>
    </QueryClientProvider>
  )
}
