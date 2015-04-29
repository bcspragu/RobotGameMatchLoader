import numpy as np

from sklearn.neighbors import KNeighborsClassifier
from sklearn.svm import SVC
from sklearn.tree import DecisionTreeClassifier

from sklearn.cross_validation import cross_val_score

#train = np.load('train.npy')[:10000]
train2 = np.load('synth.npy')[:10000]
# Remove the labels
#val = np.load('validation.npy')

data = train2[:,1:]
target = train2[:,0]

#clf = DecisionTreeClassifier()
clf = SVC()
scores = cross_val_score(clf, data, target, cv=5)
print("Accuracy: %0.2f (+/- %0.2f)" % (scores.mean(), scores.std() * 2))
