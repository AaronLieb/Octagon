import React from 'react';
import { Routes, Route, Link, useLocation } from 'react-router-dom';
import Attendees from './components/Attendees';
import Seeds from './components/Seeds';
import Conflicts from './components/Conflicts';
import { Report } from './components/report/Report';

function App() {
  const location = useLocation();

  return (
    <div>
      <nav style={{ 
        padding: '16px', 
        borderBottom: '1px solid #ccc', 
        backgroundColor: 'white',
        marginBottom: '0'
      }}>
        <Link 
          to="/attendees"
          className={location.pathname === '/attendees' ? 'button' : 'button-secondary'}
          style={{ marginRight: '8px', textDecoration: 'none' }}
        >
          Attendees
        </Link>
        <Link 
          to="/seeds"
          className={location.pathname === '/seeds' ? 'button' : 'button-secondary'}
          style={{ marginRight: '8px', textDecoration: 'none' }}
        >
          Seeds
        </Link>
        <Link 
          to="/conflicts"
          className={location.pathname === '/conflicts' ? 'button' : 'button-secondary'}
          style={{ marginRight: '8px', textDecoration: 'none' }}
        >
          Conflicts
        </Link>
        <Link 
          to="/report"
          className={location.pathname === '/report' ? 'button' : 'button-secondary'}
          style={{ textDecoration: 'none' }}
        >
          Report
        </Link>
      </nav>
      
      <Routes>
        <Route path="/" element={<Report />} />
        <Route path="/attendees" element={<Attendees />} />
        <Route path="/seeds" element={<Seeds />} />
        <Route path="/conflicts" element={<Conflicts />} />
        <Route path="/report" element={<Report />} />
      </Routes>
    </div>
  );
}

export default App;
