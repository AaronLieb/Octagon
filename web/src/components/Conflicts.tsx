import React, { useState, useEffect } from 'react';
import API_URL, { getAuthHeaders } from '../config';

interface Conflict {
  player1: string;
  player2: string;
  reason: string;
  priority: number;
  expiration: string | null;
}

interface ConflictsResponse {
  conflicts: Conflict[];
}

const Conflicts: React.FC = () => {
  const [data, setData] = useState<ConflictsResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [newConflict, setNewConflict] = useState({ 
    player1: '', 
    player2: '', 
    reason: '', 
    priority: 1, 
    expiration: '' 
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

  useEffect(() => {
    fetchConflicts();
  }, []);

  const deleteConflict = async (index: number) => {
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
    
    try {
      const response = await fetch(`${API_URL}/api/conflicts`, {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify({
          player1: newConflict.player1,
          player2: newConflict.player2,
          reason: newConflict.reason,
          priority: newConflict.priority,
          expiration: newConflict.expiration || null,
        }),
      });
      
      if (!response.ok) {
        throw new Error('Failed to add conflict');
      }
      
      setNewConflict({ player1: '', player2: '', reason: '', priority: 1, expiration: '' });
      fetchConflicts();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    }
  };

  if (loading) {
    return (
      <div className="loading">
        Loading conflicts...
      </div>
    );
  }

  return (
    <div className="container">
      <h1 className="title">Tournament Conflicts</h1>
      
      <button onClick={fetchConflicts} className="button" style={{ marginBottom: '20px' }}>
        Refresh Conflicts
      </button>

      {error && (
        <div className="error">
          Error: {error}
        </div>
      )}

      {data && (
        <>
          <div className="tournament-info">
            <h2 className="tournament-name">Conflicts from File</h2>
            <p className="attendee-count">{data.conflicts.length} conflicts found</p>
          </div>

          <div className="add-conflict" style={{ marginBottom: '20px', padding: '16px', border: '1px solid #ddd', borderRadius: '8px' }}>
            <h3>Add New Conflict</h3>
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '8px', marginBottom: '8px' }}>
              <input
                type="text"
                value={newConflict.player1}
                onChange={(e) => setNewConflict({ ...newConflict, player1: e.target.value })}
                placeholder="Player 1"
                className="input"
              />
              <input
                type="text"
                value={newConflict.player2}
                onChange={(e) => setNewConflict({ ...newConflict, player2: e.target.value })}
                placeholder="Player 2"
                className="input"
              />
            </div>
            <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr 1fr', gap: '8px', marginBottom: '8px' }}>
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
              <input
                type="datetime-local"
                value={newConflict.expiration}
                onChange={(e) => setNewConflict({ ...newConflict, expiration: e.target.value })}
                className="input"
                title="Leave empty for no expiration"
              />
            </div>
            <button onClick={addConflict} className="button">
              Add Conflict
            </button>
          </div>

          <div className="table-container">
            <table className="table">
              <thead>
                <tr>
                  <th>Player 1</th>
                  <th>Player 2</th>
                  <th>Reason</th>
                  <th>Priority</th>
                  <th>Expiration</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {data.conflicts.map((conflict, index) => (
                  <tr key={index}>
                    <td className="name">{conflict.player1}</td>
                    <td className="name">{conflict.player2}</td>
                    <td>{conflict.reason}</td>
                    <td>P{conflict.priority}</td>
                    <td>{conflict.expiration || 'Never'}</td>
                    <td>
                      <button 
                        onClick={() => deleteConflict(index)} 
                        className="button-secondary"
                        style={{ fontSize: '12px', padding: '4px 8px' }}
                      >
                        Delete
                      </button>
                    </td>
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

export default Conflicts;
