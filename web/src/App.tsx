import { Routes, Route, Navigate } from 'react-router-dom';
import Attendees from './components/Attendees';
import Seeds from './components/Seeds';
import Conflicts from './components/Conflicts';
import { Report } from './components/report/Report';
import { Login } from './components/Login';
import { NavBar } from './components/NavBar';

function App() {
  return (
    <Routes>
      <Route path="*" element={
        <>
          <NavBar />
          <Routes>
            <Route path="/" element={<Navigate to="/report" replace />} />
            <Route path="/attendees" element={<Attendees />} />
            <Route path="/seeds" element={<Seeds />} />
            <Route path="/conflicts" element={<Conflicts />} />
            <Route path="/report" element={<Report />} />
            <Route path="/login" element={<Login />} />
          </Routes>
        </>
      } />
    </Routes>
  );
}

export default App;
