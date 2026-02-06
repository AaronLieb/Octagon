import React, { useState, useEffect } from 'react';
import { useTournament } from '../utils/tournament';
import API_URL, { getAuthHeaders } from '../config';
import { DataTable } from './common/DataTable';
import { LoadingSpinner, ErrorMessage } from './common/LoadingError';
import { Attendee } from '../types';

const Attendees: React.FC = () => {
  const [attendees, setAttendees] = useState<Attendee[]>([]);
  const [tournament] = useTournament();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchAttendees = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`${API_URL}/api/attendees?tournament=${tournament}`, {
        headers: getAuthHeaders()
      });
      if (!response.ok) throw new Error('Failed to fetch attendees');
      const data = await response.json();
      setAttendees(data.attendees || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchAttendees();
  }, [tournament]);

  if (loading) return <LoadingSpinner />;

  return (
    <div className="page-container">
      <h1>Attendees ({attendees.length})</h1>
      {error && <ErrorMessage message={error} />}
      <DataTable
        data={attendees}
        columns={[
          { key: 'gamerTag', label: 'Gamer Tag' },
          { 
            key: 'name', 
            label: 'Name',
            render: (a) => `${a.firstName} ${a.lastName}`
          }
        ]}
        keyExtractor={(a) => a.id}
      />
    </div>
  );
};

export default Attendees;
