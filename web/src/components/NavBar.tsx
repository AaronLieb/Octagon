import React, { useState, useEffect } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { setTournament, getTournament } from '../utils/tournament';

export function NavBar() {
  const location = useLocation();
  const [tournamentInput, setTournamentInput] = useState(getTournament());

  useEffect(() => {
    setTournamentInput(getTournament());
  }, []);

  const handleTournamentSave = () => {
    if (tournamentInput && tournamentInput !== getTournament()) {
      setTournament(tournamentInput);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleTournamentSave();
    }
  };

  const isActive = (path: string) => location.pathname === path || (path === '/report' && location.pathname === '/');

  return (
    <nav className="navbar">
      <div className="nav-container">
        <img src="/octagon-logo.png" alt="Octagon" className="nav-logo" />
        <div className="nav-links">
          <Link to="/attendees" className={isActive('/attendees') ? 'btn btn-primary' : 'btn btn-secondary'}>
            Attendees
          </Link>
          <Link to="/seeds" className={isActive('/seeds') ? 'btn btn-primary' : 'btn btn-secondary'}>
            Seeds
          </Link>
          <Link to="/conflicts" className={isActive('/conflicts') ? 'btn btn-primary' : 'btn btn-secondary'}>
            Conflicts
          </Link>
          <Link to="/report" className={isActive('/report') ? 'btn btn-primary' : 'btn btn-secondary'}>
            Report
          </Link>
        </div>
      </div>
      <div className="tournament-selector">
        <label>Tournament:</label>
        <input
          type="text"
          value={tournamentInput}
          onChange={e => setTournamentInput(e.target.value)}
          onBlur={handleTournamentSave}
          onKeyDown={handleKeyDown}
          placeholder="octagon"
        />
      </div>
    </nav>
  );
}
