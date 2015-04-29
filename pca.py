import numpy as np
from sklearn.decomposition import PCA

# Load in the text file produced by locations.go, and run PCA on each round to
# reduce the dimensionality

filename = "rounds.csv"
data = np.loadtxt(filename, dtype=int)

# Our new representation is the number of data cases by our new basis of
# 10-elements per round, plus the label
ndata = data[:,0]
ndata = ndata.reshape((len(ndata), 1))
for i in range(100):
    print("PCA on round", i)
    dr = PCA(n_components=10)
    # Fit to one round of positions
    ndata = np.hstack((ndata, dr.fit_transform(data[:,(1 + 289 * i):(1 + 289 * (i + 1))])))

np.save(filename.split(".")[0] + ".npy", ndata)
