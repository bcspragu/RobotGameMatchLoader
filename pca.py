import numpy as np
from skimage.io import imsave
from skimage.transform import rescale
from sklearn.decomposition import PCA

data = np.loadtxt("rounds.csv", dtype=int)
# Our new representation is the number of data cases by our new basis of
# 10-elements per round, plus the label
ndata = np.zeros([len(data), 10 * 100 + 1)])
for i in range(100):
    dr = PCA(n_components=10)
    # Fit to one round of positions
    dr.fit_transform(data[:,(1 + 289 * i):(1 + 289 * (i + 1))])
