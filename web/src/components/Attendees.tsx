import React, { useState, useEffect } from 'react';
import API_URL, { getAuthHeaders } from '../config';

interface Attendee {
  id: string;
  gamerTag: string;
  firstName: string;
  lastName: string;
  playerId: string;
}

interface AttendeesResponse {
  tournament: string;
  attendees: Attendee[];
}

const Attendees: React.FC = () => {
  const [data, setData] = useState<AttendeesResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [tournament, setTournament] = useState('octagon');

  const fetchAttendees = async (tournamentSlug: string) => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch(`${API_URL}/api/attendees?tournament=${tournamentSlug}`, {
        headers: getAuthHeaders()
      });
      if (!response.ok) {
        throw new Error('Failed to fetch attendees');
      }
      const result = await response.json();
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchAttendees('octagon');
  }, []);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    fetchAttendees(tournament);
  };

  if (loading) {
    return (
      <div className="loading">
        Loading attendees...
      </div>
    );
  }

  if (error) {
    return (
      <div className="container">
        <div className="error">
          Error: {error}
        </div>
      </div>
    );
  }

  return (
    <div className="container">
      <h1 className="title">Tournament Attendees</h1>
      
      <form onSubmit={handleSubmit} className="form">
        <input
          type="text"
          value={tournament}
          onChange={(e) => setTournament(e.target.value)}
          placeholder="Tournament slug (e.g., octagon)"
          className="input"
        />
        <button type="submit" className="button">
          Load
        </button>
      </form>

      {data && (
        <>
          <div className="tournament-info">
            <h2 className="tournament-name">{data.tournament}</h2>
            <p className="attendee-count">{data.attendees.length} attendees</p>
          </div>

          <div className="table-container">
            <table className="table">
              <thead>
                <tr>
                  <th>Gamer Tag</th>
                  <th>Name</th>
                  <th>Player ID</th>
                </tr>
              </thead>
              <tbody>
                {data.attendees.map((attendee) => (
                  <tr key={attendee.id}>
                    <td className="gamer-tag">{attendee.gamerTag}</td>
                    <td className="name">{attendee.firstName} {attendee.lastName}</td>
                    <td className="player-id">{attendee.playerId}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </>
      )}
    </div>
  );
};

export default Attendees;
