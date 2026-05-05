import React, { useState, useEffect } from 'react';
import { useTournament } from '../utils/tournament';
import API_URL, { getAuthHeaders } from '../config';
import { Button } from './common/Button';
import { DataTable } from './common/DataTable';
import { LoadingSpinner, ErrorMessage } from './common/LoadingError';
import { Attendee, Conflict, ConflictsResponse } from '../types';

const Conflicts: React.FC = () => {
  const [data, setData] = useState<ConflictsResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [attendees, setAttendees] = useState<Attendee[]>([]);
  const [tournament] = useTournament();
  const [newConflict, setNewConflict] = useState({
    player1: '',
    player2: '',
    reason: '',
    priority: 1,
    expiration: '',
    round: ''
  });

  const fetchConflicts = async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(`${API_URL}/api/conflicts`, {
        headers: getAuthHeaders()
      });
      if (!response.ok) {
        throw new Error('Failed to fetch conflicts');
      }
      const result = await response.json();
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  const fetchAttendees = async () => {
    try {
      const response = await fetch(`${API_URL}/api/attendees?tournament=${tournament}`, {
        headers: getAuthHeaders()
      });
      if (!response.ok) return;
      const data = await response.json();
      setAttendees(data.attendees || []);
    } catch {
      // Silently fail — attendees are optional enhancement
    }
  };

  useEffect(() => {
    fetchConflicts();
    fetchAttendees();
  }, [tournament]);

  const deleteConflict = async (index: number) => {
    setError(null);
    try {
      const response = await fetch(`${API_URL}/api/conflicts/${index}`, {
        method: 'DELETE',
        headers: getAuthHeaders(),
      });

      if (!response.ok) {
        throw new Error('Failed to delete conflict');
      }

      fetchConflicts();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    }
  };

  const addConflict = async () => {
    if (!newConflict.player1 || !newConflict.player2) return;

    const p1 = attendees.find(a => a.playerId === newConflict.player1);
    const p2 = attendees.find(a => a.playerId === newConflict.player2);

    setError(null);
    try {
      const response = await fetch(`${API_URL}/api/conflicts`, {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify({
          player1: p1?.gamerTag || '',
          player2: p2?.gamerTag || '',
          player1Id: newConflict.player1,
          player2Id: newConflict.player2,
          reason: newConflict.reason,
          priority: newConflict.priority,
          expiration: newConflict.expiration || null,
          round: newConflict.round || null,
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to add conflict');
      }

      setNewConflict({ player1: '', player2: '', reason: '', priority: 1, expiration: '', round: '' });
      fetchConflicts();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    }
  };

  if (loading) return <LoadingSpinner />;

  const sortedAttendees = [...attendees].sort((a, b) => a.gamerTag.localeCompare(b.gamerTag));

  return (
    <div className="page-container">
      <h1>Conflicts</h1>

      <Button onClick={fetchConflicts} style={{ marginBottom: '20px' }}>
        Refresh Conflicts
      </Button>

      {error && <ErrorMessage message={error} />}

      {data && (
        <>
          <p style={{ color: '#666', marginBottom: '20px' }}>{data.conflicts.length} conflicts found</p>

          <div style={{ marginBottom: '20px', padding: '16px', border: `1px solid var(--border-color)`, borderRadius: '8px', background: 'var(--bg-secondary)' }}>
            <h2>Add New Conflict</h2>
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '8px', marginBottom: '8px' }}>
              <select
                value={newConflict.player1}
                onChange={(e) => setNewConflict({ ...newConflict, player1: e.target.value })}
                className="input"
              >
                <option value="">Select Player 1</option>
                {sortedAttendees.map(a => (
                  <option key={a.playerId} value={a.playerId}>{a.gamerTag}</option>
                ))}
              </select>
              <select
                value={newConflict.player2}
                onChange={(e) => setNewConflict({ ...newConflict, player2: e.target.value })}
                className="input"
              >
                <option value="">Select Player 2</option>
                {sortedAttendees.map(a => (
                  <option key={a.playerId} value={a.playerId}>{a.gamerTag}</option>
                ))}
              </select>
            </div>
            <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr 1fr 1fr', gap: '8px', marginBottom: '8px' }}>
              <input
                type="text"
                value={newConflict.reason}
                onChange={(e) => setNewConflict({ ...newConflict, reason: e.target.value })}
                placeholder="Reason"
                className="input"
              />
              <div>
                <label style={{ fontSize: '12px', color: '#666' }}>Priority</label>
                <input
                  type="number"
                  value={newConflict.priority}
                  onChange={(e) => setNewConflict({ ...newConflict, priority: parseInt(e.target.value) || 1 })}
                  placeholder="Priority"
                  className="input"
                  min="1"
                />
              </div>
              <div>
                <label style={{ fontSize: '12px', color: '#666' }}>Round</label>
                <input
                  type="text"
                  value={newConflict.round}
                  onChange={(e) => setNewConflict({ ...newConflict, round: e.target.value })}
                  placeholder="e.g. WR1, LR2"
                  className="input"
                />
              </div>
              <input
                type="datetime-local"
                value={newConflict.expiration}
                onChange={(e) => setNewConflict({ ...newConflict, expiration: e.target.value })}
                className="input"
                title="Leave empty for no expiration"
              />
            </div>
            <Button onClick={addConflict}>Add Conflict</Button>
          </div>

          <DataTable
            data={data.conflicts}
            columns={[
              { key: 'player1', label: 'Player 1' },
              { key: 'player2', label: 'Player 2' },
              { key: 'reason', label: 'Reason' },
              { 
                key: 'priority', 
                label: 'Priority',
                render: (c) => c.priority < 0 ? `R${-c.priority}` : `P${c.priority}`
              },
              {
                key: 'round',
                label: 'Round',
                render: (c) => c.round || 'All'
              },
              { 
                key: 'expiration', 
                label: 'Expiration',
                render: (c) => c.expiration || 'Never'
              },
              {
                key: 'actions',
                label: 'Actions',
                render: (_, idx) => (
                  <Button 
                    variant="secondary" 
                    onClick={() => deleteConflict(idx)}
                    style={{ fontSize: '12px', padding: '4px 8px' }}
                  >
                    Delete
                  </Button>
                )
              }
            ]}
            keyExtractor={(_, idx) => idx}
          />
        </>
      )}
    </div>
  );
};

export default Conflicts;
