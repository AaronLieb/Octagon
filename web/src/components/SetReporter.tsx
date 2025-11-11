import React, { useState, useEffect } from 'react';

interface Player {
  name: string;
  id: number;
}

interface Set {
  id: number;
  player1: Player;
  player2: Player;
  round: string;
  entrant1: number;
  entrant2: number;
}

interface GameResult {
  winner: number;
  p1Char: string;
  p2Char: string;
}

interface SetsResponse {
  sets: Set[];
}

const SetReporter: React.FC = () => {
  const [sets, setSets] = useState<Set[]>([]);
  const [filteredSets, setFilteredSets] = useState<Set[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [tournament, setTournament] = useState('octagon');
  const [redemption, setRedemption] = useState(false);
  const [filter, setFilter] = useState('');
  const [selectedSet, setSelectedSet] = useState<Set | null>(null);
  const [games, setGames] = useState<GameResult[]>([
    { winner: 0, p1Char: '', p2Char: '' },
    { winner: 0, p1Char: '', p2Char: '' },
    { winner: 0, p1Char: '', p2Char: '' },
    { winner: 0, p1Char: '', p2Char: '' },
    { winner: 0, p1Char: '', p2Char: '' },
  ]);

  const fetchSets = async () => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch(`http://localhost:8080/api/sets?tournament=${tournament}&redemption=${redemption}`);
      if (!response.ok) {
        throw new Error('Failed to fetch sets');
      }
      const result: SetsResponse = await response.json();
      setSets(result.sets);
      setFilteredSets(result.sets);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (filter === '') {
      setFilteredSets(sets);
    } else {
      const filtered = sets.filter(set => 
        set.player1.name.toLowerCase().includes(filter.toLowerCase()) ||
        set.player2.name.toLowerCase().includes(filter.toLowerCase())
      );
      setFilteredSets(filtered);
    }
  }, [filter, sets]);

  const selectSet = (set: Set) => {
    setSelectedSet(set);
    setGames([
      { winner: 0, p1Char: '', p2Char: '' },
      { winner: 0, p1Char: '', p2Char: '' },
      { winner: 0, p1Char: '', p2Char: '' },
      { winner: 0, p1Char: '', p2Char: '' },
      { winner: 0, p1Char: '', p2Char: '' },
    ]);
  };

  const updateGame = (gameIndex: number, field: keyof GameResult, value: string | number) => {
    const newGames = [...games];
    newGames[gameIndex] = { ...newGames[gameIndex], [field]: value };
    
    // Copy characters from previous game if not set
    if (gameIndex > 0 && field === 'winner') {
      if (!newGames[gameIndex].p1Char) {
        newGames[gameIndex].p1Char = newGames[gameIndex - 1].p1Char;
      }
      if (!newGames[gameIndex].p2Char) {
        newGames[gameIndex].p2Char = newGames[gameIndex - 1].p2Char;
      }
    }
    
    setGames(newGames);
  };

  const reportSet = async () => {
    if (!selectedSet) return;
    
    const playedGames = games.filter(game => game.winner > 0);
    if (playedGames.length === 0) return;
    
    try {
      const response = await fetch('http://localhost:8080/api/sets/report', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          setId: selectedSet.id,
          games: playedGames,
        }),
      });
      
      if (!response.ok) {
        throw new Error('Failed to report set');
      }
      
      setSelectedSet(null);
      fetchSets(); // Refresh sets
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      reportSet();
    }
  };

  if (selectedSet) {
    return (
      <div className="container" onKeyPress={handleKeyPress} tabIndex={0}>
        <h1 className="title">Report Set</h1>
        
        <div style={{ marginBottom: '20px', padding: '16px', border: '1px solid #ddd', borderRadius: '8px' }}>
          <h2>{selectedSet.player1.name} vs {selectedSet.player2.name}</h2>
          <p>Round: {selectedSet.round}</p>
          <button onClick={() => setSelectedSet(null)} className="button-secondary">
            Back to Sets
          </button>
        </div>

        <div style={{ marginBottom: '20px' }}>
          <h3>Games (5 max)</h3>
          {games.map((game, index) => (
            <div key={index} style={{ 
              display: 'grid', 
              gridTemplateColumns: '50px 1fr 1fr 1fr 1fr', 
              gap: '8px', 
              marginBottom: '8px',
              alignItems: 'center'
            }}>
              <span>G{index + 1}:</span>
              <select
                value={game.winner}
                onChange={(e) => updateGame(index, 'winner', parseInt(e.target.value))}
                className="input"
              >
                <option value={0}>-</option>
                <option value={1}>{selectedSet.player1.name}</option>
                <option value={2}>{selectedSet.player2.name}</option>
              </select>
              <input
                type="text"
                value={game.p1Char}
                onChange={(e) => updateGame(index, 'p1Char', e.target.value)}
                placeholder={`${selectedSet.player1.name} character`}
                className="input"
              />
              <input
                type="text"
                value={game.p2Char}
                onChange={(e) => updateGame(index, 'p2Char', e.target.value)}
                placeholder={`${selectedSet.player2.name} character`}
                className="input"
              />
              <span>{game.winner > 0 ? '✓' : ''}</span>
            </div>
          ))}
        </div>

        <button onClick={reportSet} className="button">
          Report Set (Enter)
        </button>
        
        {error && (
          <div className="error" style={{ marginTop: '10px' }}>
            Error: {error}
          </div>
        )}
      </div>
    );
  }

  return (
    <div className="container">
      <h1 className="title">Set Reporter</h1>
      
      <div className="form" style={{ marginBottom: '20px' }}>
        <div style={{ display: 'flex', gap: '8px', marginBottom: '8px' }}>
          <input
            type="text"
            value={tournament}
            onChange={(e) => setTournament(e.target.value)}
            placeholder="Tournament slug"
            className="input"
          />
          <label style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
            <input
              type="checkbox"
              checked={redemption}
              onChange={(e) => setRedemption(e.target.checked)}
            />
            Redemption
          </label>
          <button onClick={fetchSets} className="button">
            Load Sets
          </button>
        </div>
        
        <input
          type="text"
          value={filter}
          onChange={(e) => setFilter(e.target.value)}
          placeholder="Filter by player name..."
          className="input"
        />
      </div>

      {loading && (
        <div className="loading">Loading sets...</div>
      )}

      {error && (
        <div className="error">Error: {error}</div>
      )}

      {filteredSets.length > 0 && (
        <div className="table-container">
          <table className="table">
            <thead>
              <tr>
                <th>Round</th>
                <th>Player 1</th>
                <th>Player 2</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {filteredSets.map((set) => (
                <tr key={set.id}>
                  <td>{set.round}</td>
                  <td className="name">{set.player1.name}</td>
                  <td className="name">{set.player2.name}</td>
                  <td>
                    <button 
                      onClick={() => selectSet(set)} 
                      className="button"
                      style={{ fontSize: '12px', padding: '4px 8px' }}
                    >
                      Report
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

export default SetReporter;
